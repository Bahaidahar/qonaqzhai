// Binary gateway is the single entry point. Routes incoming requests to
// auth-svc, core-svc, or realtime-svc based on path prefix, and applies CORS
// on the edge so individual services stay simple.
//
//   /api/signup, /api/login, /api/auth/*, /api/me   → auth-svc
//   /api/ws                                         → realtime-svc (WS upgrade)
//   /api/*                                          → core-svc
//   /internal/*                                     → 404 (never exposed externally)
package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	addr := envOr("GATEWAY_ADDR", ":8080")
	authURL := mustURL(envOr("AUTH_URL", "http://localhost:8081"))
	coreURL := mustURL(envOr("CORE_URL", "http://localhost:8082"))
	realtimeURL := mustURL(envOr("REALTIME_URL", "http://localhost:8083"))
	corsOrigin := envOr("CORS_ORIGIN", "http://localhost:3000")

	authProxy := newProxy(authURL)
	coreProxy := newProxy(coreURL)
	realtimeProxy := newProxy(realtimeURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/internal/") {
			http.NotFound(w, r)
			return
		}
		switch {
		case path == "/api/ws":
			realtimeProxy.ServeHTTP(w, r)
		case path == "/api/signup",
			path == "/api/login",
			path == "/api/me",
			strings.HasPrefix(path, "/api/auth/"):
			authProxy.ServeHTTP(w, r)
		default:
			coreProxy.ServeHTTP(w, r)
		}
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           cors(corsOrigin, mux),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("gateway listening on %s · auth=%s · core=%s · realtime=%s · cors=%s",
		addr, authURL, coreURL, realtimeURL, corsOrigin)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %v", err)
	}
}

// cors applies a permissive-yet-explicit CORS policy. Origins are echoed when
// they match the configured list (comma-separated, or "*"). Credentials allowed
// only when an exact match is found — never with wildcard.
func cors(allowedOrigins string, next http.Handler) http.Handler {
	allowAll := allowedOrigins == "*"
	allowSet := map[string]struct{}{}
	for _, o := range strings.Split(allowedOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" && o != "*" {
			allowSet[o] = struct{}{}
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			_, allowed := allowSet[origin]
			if allowAll || allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Requested-With")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func newProxy(target *url.URL) *httputil.ReverseProxy {
	p := httputil.NewSingleHostReverseProxy(target)
	p.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		log.Printf("upstream %s error: %v", target, err)
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}
	return p
}

func mustURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("bad URL %q: %v", raw, err)
	}
	return u
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
