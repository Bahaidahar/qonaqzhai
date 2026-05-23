// Package proxy implements the reverse-proxy routing rules used by the
// gateway. Each Route picks the upstream by URL prefix; the chosen upstream is
// reached via its base URL (host + port) read from the environment in main.
package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/httpx"
)

// Route maps an incoming URL prefix to an upstream base URL.
type Route struct {
	Prefix string // matched against r.URL.Path
	Target string // upstream base URL, e.g. "http://auth:8081"
}

// Router selects + forwards requests to the appropriate upstream.
type Router struct {
	routes []routeProxy
}

type routeProxy struct {
	Route
	proxy *httputil.ReverseProxy
}

// New builds a router. Routes are matched in declaration order; the first
// matching prefix wins, so list more specific prefixes before catch-alls.
func New(routes []Route) (*Router, error) {
	rps := make([]routeProxy, 0, len(routes))
	for _, r := range routes {
		u, err := url.Parse(r.Target)
		if err != nil {
			return nil, err
		}
		rp := httputil.NewSingleHostReverseProxy(u)
		base := rp.Director
		host := u.Host
		rp.Director = func(req *http.Request) {
			base(req)
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Host = host
		}
		rps = append(rps, routeProxy{Route: r, proxy: rp})
	}
	return &Router{routes: rps}, nil
}

// ServeHTTP picks the first prefix match and forwards. 404 when nothing matches.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, rp := range r.routes {
		if strings.HasPrefix(req.URL.Path, rp.Prefix) {
			rp.proxy.ServeHTTP(w, req)
			return
		}
	}
	httpx.WriteError(w, http.StatusNotFound, "no route")
}

// WithAuthForwarding wraps next so that authenticated requests get the
// canonical X-User-Id / X-User-Role / X-User-Email headers stamped before
// proxying. Anonymous requests have any forged values stripped.
func WithAuthForwarding(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c, ok := pkgauth.FromContext(r.Context()); ok {
			r.Header.Set("X-User-Id", c.UserID)
			r.Header.Set("X-User-Role", c.Role)
			r.Header.Set("X-User-Email", c.Email)
		} else {
			r.Header.Del("X-User-Id")
			r.Header.Del("X-User-Role")
			r.Header.Del("X-User-Email")
		}
		next.ServeHTTP(w, r)
	})
}
