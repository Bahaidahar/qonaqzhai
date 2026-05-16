package pay_test

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"qonaqzhai-backend/internal/adapter/pay"
	"qonaqzhai-backend/internal/usecase"
)

func TestNewDisabledWithoutCredentials(t *testing.T) {
	t.Parallel()
	if pay.New(pay.Config{}) != nil {
		t.Error("expected nil when credentials missing")
	}
}

func TestCreatePaymentSendsSignedRequest(t *testing.T) {
	t.Parallel()
	var seenBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		seenBody = string(b)
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(`<?xml version="1.0"?>
		<response>
		  <pg_status>ok</pg_status>
		  <pg_payment_id>pi_12345</pg_payment_id>
		  <pg_redirect_url>https://paybox.example/checkout/pi_12345</pg_redirect_url>
		</response>`))
	}))
	defer srv.Close()

	gw := pay.New(pay.Config{
		MerchantID: "M1", SecretKey: "secret", Sandbox: true, BaseURL: srv.URL, Client: srv.Client(),
	})
	if gw == nil {
		t.Fatal("expected gateway")
	}
	out, err := gw.CreatePayment(context.Background(), usecase.PaymentIntent{
		OrderID: "b-1", Amount: 1500, Currency: "KZT", CustomerEmail: "x@y.kz",
		SuccessURL: "https://app/success", FailureURL: "https://app/fail",
		Description: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.TransactionID != "pi_12345" {
		t.Errorf("transactionId=%s", out.TransactionID)
	}
	if out.URL == "" {
		t.Error("missing redirect URL")
	}
	if !strings.Contains(seenBody, "pg_merchant_id=M1") {
		t.Errorf("merchant_id missing from form: %s", seenBody)
	}
	if !strings.Contains(seenBody, "pg_sig=") {
		t.Error("signature missing from form")
	}
}

func TestVerifyCallback(t *testing.T) {
	t.Parallel()
	gw := pay.New(pay.Config{MerchantID: "M1", SecretKey: "secret"})
	form := map[string]string{
		"pg_order_id":   "b-1",
		"pg_payment_id": "pi_1",
		"pg_amount":     "1500",
		"pg_result":     "1",
	}
	form["pg_sig"] = signForTest("result", form, "secret")
	res, err := gw.VerifyCallback(form)
	if err != nil {
		t.Fatal(err)
	}
	if res.OrderID != "b-1" || res.TransactionID != "pi_1" || res.Amount != 1500 || !res.Success {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestVerifyCallbackRejectsBadSignature(t *testing.T) {
	t.Parallel()
	gw := pay.New(pay.Config{MerchantID: "M1", SecretKey: "secret"})
	form := map[string]string{
		"pg_order_id": "b-1", "pg_payment_id": "pi_1", "pg_amount": "1500",
		"pg_result": "1", "pg_sig": "deadbeef",
	}
	if _, err := gw.VerifyCallback(form); err == nil {
		t.Error("bad signature accepted")
	}
}

func signForTest(script string, params map[string]string, secret string) string {
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
