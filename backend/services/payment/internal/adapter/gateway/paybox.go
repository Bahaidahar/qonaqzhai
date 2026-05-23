package gateway

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

	"qonaqzhai-backend/services/payment/internal/ports"
)

// PayBoxConfig bundles PayBox merchant credentials.
type PayBoxConfig struct {
	MerchantID string
	SecretKey  string
	Sandbox    bool
	BaseURL    string
	Client     *http.Client
}

// PayBox is the Freedom Pay (PayBox.money) PSP adapter. It currently runs the
// init_payment.php flow synchronously and returns the redirect URL as
// providerRef — the diploma's Charge path treats that as success because we
// already trust the customer's saved card. Real production code should split
// Charge into Init + capture callback.
type PayBox struct{ cfg PayBoxConfig }

// NewPayBox returns nil when credentials are missing.
func NewPayBox(cfg PayBoxConfig) *PayBox {
	if cfg.MerchantID == "" || cfg.SecretKey == "" {
		return nil
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.paybox.money"
	}
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 20 * time.Second}
	}
	return &PayBox{cfg: cfg}
}

// Charge runs init_payment.php and returns the PayBox payment id.
func (p *PayBox) Charge(ctx context.Context, in ports.ChargeInput) (string, error) {
	if p == nil {
		return "", errors.New("paybox disabled")
	}
	params := map[string]string{
		"pg_merchant_id":  p.cfg.MerchantID,
		"pg_order_id":     in.OrderID,
		"pg_amount":       strconv.FormatInt(in.Amount, 10),
		"pg_currency":     currencyOrDefault(in.Currency),
		"pg_description":  "Qonaqzhai " + in.OrderID,
		"pg_salt":         salt(),
		"pg_testing_mode": boolToStr(p.cfg.Sandbox),
	}
	params["pg_sig"] = sign("init_payment.php", params, p.cfg.SecretKey)

	form := url.Values{}
	for k, v := range params {
		if v != "" {
			form.Set(k, v)
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		p.cfg.BaseURL+"/init_payment.php", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := p.cfg.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("paybox init: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("paybox status %d: %s", res.StatusCode, string(body))
	}
	parsed := parseInit(string(body))
	if parsed.status != "ok" {
		return "", fmt.Errorf("paybox declined: %s %s", parsed.status, parsed.errorDesc)
	}
	return parsed.paymentID, nil
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

func currencyOrDefault(c string) string {
	if c == "" {
		return "KZT"
	}
	return c
}

type initResp struct{ status, paymentID, errorDesc string }

func parseInit(xml string) initResp {
	return initResp{
		status:    between(xml, "<pg_status>", "</pg_status>"),
		paymentID: between(xml, "<pg_payment_id>", "</pg_payment_id>"),
		errorDesc: between(xml, "<pg_error_description>", "</pg_error_description>"),
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
