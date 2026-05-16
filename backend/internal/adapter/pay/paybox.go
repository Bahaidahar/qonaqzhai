// Package pay implements the PaymentGateway port against Freedom Pay (PayBox.money).
//
// API reference: https://docs.paybox.money/
// The gateway uses URL-encoded forms and a SHA1 signature over sorted parameters.
package pay

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"qonaqzhai-backend/internal/usecase"
)

// Config bundles PayBox merchant credentials.
type Config struct {
	MerchantID string
	SecretKey  string
	Sandbox    bool
	BaseURL    string // overridden in tests; defaults are derived from Sandbox.
	Client     *http.Client
}

// PayBox implements usecase.PaymentGateway.
type PayBox struct {
	cfg Config
}

// New returns nil when credentials are not supplied (payment disabled).
func New(cfg Config) *PayBox {
	if cfg.MerchantID == "" || cfg.SecretKey == "" {
		return nil
	}
	if cfg.BaseURL == "" {
		if cfg.Sandbox {
			cfg.BaseURL = "https://api.paybox.money"
		} else {
			cfg.BaseURL = "https://api.paybox.money"
		}
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 20 * time.Second}
	}
	return &PayBox{cfg: cfg}
}

// CreatePayment initialises a payment and returns the gateway redirect URL.
//
// PayBox flow:
//   POST /init_payment.php  →  pg_payment_id + pg_redirect_url
//
// All params are URL-encoded; signature is SHA1 over the sorted values
// joined with `;` plus the merchant secret key.
func (p *PayBox) CreatePayment(ctx context.Context, in usecase.PaymentIntent) (usecase.PaymentRedirect, error) {
	if p == nil {
		return usecase.PaymentRedirect{}, errors.New("paybox disabled")
	}
	params := map[string]string{
		"pg_merchant_id":     p.cfg.MerchantID,
		"pg_order_id":        in.OrderID,
		"pg_amount":          strconv.FormatInt(in.Amount, 10),
		"pg_currency":        defaultCurrency(in.Currency),
		"pg_description":     defaultDescription(in.Description, in.OrderID),
		"pg_salt":            salt(),
		"pg_user_contact_email": in.CustomerEmail,
		"pg_success_url":     in.SuccessURL,
		"pg_failure_url":     in.FailureURL,
		"pg_result_url":      strings.TrimRight(in.SuccessURL, "/") + "/callback",
		"pg_testing_mode":    boolToStr(p.cfg.Sandbox),
	}
	params["pg_sig"] = sign("init_payment.php", params, p.cfg.SecretKey)

	form := url.Values{}
	for k, v := range params {
		if v == "" {
			continue
		}
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.cfg.BaseURL+"/init_payment.php", strings.NewReader(form.Encode()))
	if err != nil {
		return usecase.PaymentRedirect{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := p.cfg.Client.Do(req)
	if err != nil {
		return usecase.PaymentRedirect{}, fmt.Errorf("paybox init: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return usecase.PaymentRedirect{}, fmt.Errorf("paybox status %d: %s", res.StatusCode, string(body))
	}
	// PayBox returns XML by default; parse minimal fields.
	parsed := parseInitResponse(string(body))
	if parsed.status != "ok" {
		return usecase.PaymentRedirect{}, fmt.Errorf("paybox declined: status=%s desc=%s", parsed.status, parsed.errorDesc)
	}
	return usecase.PaymentRedirect{
		URL:           parsed.redirectURL,
		TransactionID: parsed.paymentID,
	}, nil
}

// VerifyCallback validates the pg_sig signature on a result-callback form.
// Caller maps form fields from the incoming request.
func (p *PayBox) VerifyCallback(form map[string]string) (usecase.CallbackResult, error) {
	if p == nil {
		return usecase.CallbackResult{}, errors.New("paybox disabled")
	}
	sig := form["pg_sig"]
	delete(form, "pg_sig")
	want := sign("result", form, p.cfg.SecretKey)
	if !strings.EqualFold(sig, want) {
		return usecase.CallbackResult{}, errors.New("invalid paybox signature")
	}
	amount, _ := strconv.ParseInt(form["pg_amount"], 10, 64)
	return usecase.CallbackResult{
		OrderID:       form["pg_order_id"],
		TransactionID: form["pg_payment_id"],
		Amount:        amount,
		Success:       strings.EqualFold(form["pg_result"], "1"),
	}, nil
}

func sign(script string, params map[string]string, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := []string{script}
	for _, k := range keys {
		parts = append(parts, params[k])
	}
	parts = append(parts, secret)
	h := sha1.Sum([]byte(strings.Join(parts, ";")))
	return hex.EncodeToString(h[:])
}

func salt() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }

func boolToStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func defaultCurrency(c string) string {
	if c == "" {
		return "KZT"
	}
	return c
}

func defaultDescription(d, fallback string) string {
	if d == "" {
		return "Qonaqzhai booking " + fallback
	}
	return d
}

type initResponse struct {
	status      string
	paymentID   string
	redirectURL string
	errorDesc   string
}

// parseInitResponse extracts fields from the XML envelope without pulling
// in an XML parser dependency — values live inside simple <tag>…</tag> pairs.
func parseInitResponse(xml string) initResponse {
	return initResponse{
		status:      between(xml, "<pg_status>", "</pg_status>"),
		paymentID:   between(xml, "<pg_payment_id>", "</pg_payment_id>"),
		redirectURL: between(xml, "<pg_redirect_url>", "</pg_redirect_url>"),
		errorDesc:   between(xml, "<pg_error_description>", "</pg_error_description>"),
	}
}

func between(s, start, end string) string {
	i := strings.Index(s, start)
	if i < 0 {
		return ""
	}
	rest := s[i+len(start):]
	j := strings.Index(rest, end)
	if j < 0 {
		return ""
	}
	return strings.TrimSpace(rest[:j])
}

var _ usecase.PaymentGateway = (*PayBox)(nil)
