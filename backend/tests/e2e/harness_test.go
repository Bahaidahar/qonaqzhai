// Package e2e brings up the entire microservice stack (4 Postgres containers
// + 5 service binaries) and exercises critical flows through the gateway.
//
// Build tag — Docker is required and bringup takes ~30s, so run explicitly:
//
//	go test -tags=e2e ./tests/e2e -v

//go:build e2e

package e2e_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// dsnFor brings up a fresh Postgres container and returns its libpq DSN.
func dsnFor(t *testing.T, ctx context.Context, dbName string) string {
	t.Helper()
	pg, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		t.Fatalf("postgres %s: %v", dbName, err)
	}
	t.Cleanup(func() { _ = pg.Terminate(context.Background()) })
	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("dsn %s: %v", dbName, err)
	}
	time.Sleep(1 * time.Second)
	return dsn
}

// startBin compiles the service binary once and runs it with the supplied env.
// Each service is killed (process group) when the test ends. Compiling
// up-front avoids `go run` spawning a child that ignores SIGKILL on its
// parent and lets cleanup return promptly.
func startBin(t *testing.T, name string, env map[string]string) {
	t.Helper()
	wd, _ := os.Getwd()
	root := filepath.Join(wd, "..", "..")
	bin := filepath.Join(t.TempDir(), name)
	buildCmd := exec.Command("go", "build", "-o", bin, "./services/"+name+"/cmd/"+name)
	buildCmd.Dir = root
	buildCmd.Env = os.Environ()
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("build %s: %v", name, err)
	}

	cmd := exec.Command(bin)
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start %s: %v", name, err)
	}
	t.Cleanup(func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		done := make(chan struct{})
		go func() { _, _ = cmd.Process.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})
}

// waitFor probes url until it returns 200 or timeout expires.
func waitFor(t *testing.T, url string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url) //nolint:gosec
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for %s", url)
}

// --- HTTP helpers ------------------------------------------------------------

func req(t *testing.T, method, url, body, token string, wantStatus int) map[string]any {
	t.Helper()
	var rd io.Reader
	if body != "" {
		rd = strings.NewReader(body)
	}
	r, err := http.NewRequest(method, url, rd)
	if err != nil {
		t.Fatalf("build req: %v", err)
	}
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s status=%d want=%d body=%s", method, url, resp.StatusCode, wantStatus, raw)
	}
	if len(raw) == 0 || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s: %v body=%s", url, err, raw)
	}
	return out
}

func reqList(t *testing.T, method, url, token string, wantStatus int) []any {
	t.Helper()
	r, _ := http.NewRequest(method, url, nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s status=%d want=%d body=%s", method, url, resp.StatusCode, wantStatus, raw)
	}
	var out []any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode list %s: %v body=%s", url, err, raw)
	}
	return out
}

// --- bootstrap ---------------------------------------------------------------

type stack struct {
	GatewayURL string
}

// bootStack spins up all 4 dbs + all 5 services and waits for them to be ready.
// Returns the gateway base URL.
func bootStack(t *testing.T) *stack {
	t.Helper()
	if testing.Short() {
		t.Skip("e2e suite")
	}
	ctx := context.Background()

	authDSN := dsnFor(t, ctx, "auth")
	coreDSN := dsnFor(t, ctx, "core")
	paymentDSN := dsnFor(t, ctx, "payment")
	realtimeDSN := dsnFor(t, ctx, "realtime")

	startBin(t, "auth", map[string]string{
		"AUTH_DATABASE_URL": authDSN,
		"AUTH_HTTP_ADDR":    ":18181",
		"AUTH_GRPC_ADDR":    ":19181",
		"JWT_SECRET":        "test-secret",
		"ADMIN_EMAIL":       "admin@qonaqzhai.kz",
		"ADMIN_PASSWORD":    "adminpass1",
	})
	startBin(t, "core", map[string]string{
		"CORE_DATABASE_URL":  coreDSN,
		"CORE_HTTP_ADDR":     ":18282",
		"CORE_GRPC_ADDR":     ":19282",
		"AUTH_GRPC_ADDR":     "localhost:19181",
		"PAYMENT_GRPC_ADDR":  "localhost:19383",
		"REALTIME_GRPC_ADDR": "localhost:19484",
	})
	startBin(t, "payment", map[string]string{
		"PAYMENT_DATABASE_URL": paymentDSN,
		"PAYMENT_HTTP_ADDR":    ":18383",
		"PAYMENT_GRPC_ADDR":    ":19383",
		"AUTH_GRPC_ADDR":       "localhost:19181",
		"CORE_GRPC_ADDR":       "localhost:19282",
	})
	startBin(t, "realtime", map[string]string{
		"REALTIME_DATABASE_URL": realtimeDSN,
		"REALTIME_HTTP_ADDR":    ":18484",
		"REALTIME_GRPC_ADDR":    ":19484",
		"AUTH_GRPC_ADDR":        "localhost:19181",
	})
	startBin(t, "gateway", map[string]string{
		"GATEWAY_ADDR":      ":18080",
		"AUTH_GRPC_ADDR":    "localhost:19181",
		"AUTH_HTTP_URL":     "http://localhost:18181",
		"CORE_HTTP_URL":     "http://localhost:18282",
		"PAYMENT_HTTP_URL":  "http://localhost:18383",
		"REALTIME_HTTP_URL": "http://localhost:18484",
	})

	// Wait for each service's health endpoint.
	waitFor(t, "http://localhost:18181/api/health", 40*time.Second)
	waitFor(t, "http://localhost:18282/api/health", 40*time.Second)
	waitFor(t, "http://localhost:18383/api/health", 40*time.Second)
	waitFor(t, "http://localhost:18484/api/health", 40*time.Second)
	// Gateway is ready when it can forward to /api/health on core.
	waitFor(t, "http://localhost:18080/api/health", 40*time.Second)

	return &stack{GatewayURL: "http://localhost:18080"}
}

// --- tests -------------------------------------------------------------------

// TestSignupLoginMe is the smoke check for the auth path through gateway.
func TestSignupLoginMe(t *testing.T) {
	s := bootStack(t)

	req(t, "POST", s.GatewayURL+"/api/signup",
		`{"email":"smoke@qonaqzhai.kz","password":"password123","name":"Smoke"}`, "",
		http.StatusCreated)

	login := req(t, "POST", s.GatewayURL+"/api/login",
		`{"email":"smoke@qonaqzhai.kz","password":"password123"}`, "",
		http.StatusOK)

	token, _ := login["accessToken"].(string)
	if token == "" {
		t.Fatal("missing access token")
	}

	me := req(t, "GET", s.GatewayURL+"/api/me", "", token, http.StatusOK)
	if me["email"] != "smoke@qonaqzhai.kz" {
		t.Fatalf("/api/me email mismatch: %v", me)
	}
}

// TestFullBookingFlow exercises the cross-service happy path:
//
//  1. customer signs up + logs in
//  2. vendor signs up, creates a vendor profile (pending)
//  3. admin approves the vendor
//  4. customer browses public catalog and finds the vendor
//  5. customer creates a booking
//  6. vendor accepts; realtime thread is auto-created via gRPC
//  7. customer + vendor exchange a message in the thread
//  8. customer adds a card on payment-svc and pays the booking
//  9. booking flips to paid via the atomic MarkPaid saga
//  10. customer leaves a 5-star review which updates the vendor rating
func TestFullBookingFlow(t *testing.T) {
	s := bootStack(t)
	gw := s.GatewayURL

	// 1. customer
	custTok := login(t, gw, "cust@x.kz", "password123", "Cust", "customer")

	// 2. vendor signup + profile
	vendTok := login(t, gw, "vend@x.kz", "password123", "Vend", "vendor")
	req(t, "PUT", gw+"/api/me/vendor",
		`{"name":"Rixos","category":"venue","city":"Almaty","description":"d","priceFrom":100000}`,
		vendTok, http.StatusOK)

	// 3. admin approves
	adminTok := mustLogin(t, gw, "admin@qonaqzhai.kz", "adminpass1")
	myVendor := req(t, "GET", gw+"/api/me/vendor", "", vendTok, http.StatusOK)
	vendorID, _ := myVendor["id"].(string)
	if vendorID == "" {
		t.Fatal("vendor id missing")
	}
	req(t, "PATCH", gw+"/api/admin/vendors/"+vendorID+"/status",
		`{"status":"approved"}`, adminTok, http.StatusOK)

	// 4. public catalog
	catalog := req(t, "GET", gw+"/api/vendors", "", "", http.StatusOK)
	items, _ := catalog["items"].([]any)
	if len(items) == 0 {
		t.Fatal("approved vendor not in catalog")
	}

	// 5. booking
	booking := req(t, "POST", gw+"/api/bookings",
		`{"vendorId":"`+vendorID+`","eventDate":"2026-09-01","guestCount":50,"amount":200000}`,
		custTok, http.StatusCreated)
	bookingID, _ := booking["id"].(string)
	if bookingID == "" {
		t.Fatal("booking id missing")
	}

	// 6. vendor accepts (also triggers realtime EnsureThread)
	req(t, "PATCH", gw+"/api/bookings/"+bookingID,
		`{"status":"accepted"}`, vendTok, http.StatusOK)

	// Realtime takes a beat to receive the gRPC.
	time.Sleep(500 * time.Millisecond)

	// 7. customer + vendor exchange messages
	threads := reqList(t, "GET", gw+"/api/threads", custTok, http.StatusOK)
	if len(threads) == 0 {
		t.Fatal("no thread created on accept")
	}
	threadID, _ := threads[0].(map[string]any)["thread"].(map[string]any)["id"].(string)
	if threadID == "" {
		t.Fatal("thread id missing")
	}
	req(t, "POST", gw+"/api/threads/"+threadID+"/messages",
		`{"text":"hi, looking forward"}`, custTok, http.StatusCreated)
	req(t, "POST", gw+"/api/threads/"+threadID+"/messages",
		`{"text":"thanks, see you on the 1st"}`, vendTok, http.StatusCreated)

	thread := req(t, "GET", gw+"/api/threads/"+threadID, "", vendTok, http.StatusOK)
	msgs, _ := thread["messages"].([]any)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}

	// 8. customer adds a card + pays
	card := req(t, "POST", gw+"/api/cards",
		`{"number":"4111111111111111","expMonth":6,"expYear":30,"holder":"TEST"}`,
		custTok, http.StatusCreated)
	cardID, _ := card["id"].(string)
	if cardID == "" {
		t.Fatal("card id missing")
	}
	paid := req(t, "POST", gw+"/api/bookings/"+bookingID+"/pay",
		`{"cardId":"`+cardID+`","currency":"KZT"}`, custTok, http.StatusOK)
	if paid["status"] != "paid" {
		t.Fatalf("expected paid, got %v", paid["status"])
	}
	if paid["paymentId"] == "" {
		t.Fatal("paymentId missing on paid booking")
	}

	// 9. customer reviews
	req(t, "POST", gw+"/api/reviews",
		`{"bookingId":"`+bookingID+`","rating":5,"text":"great"}`,
		custTok, http.StatusCreated)

	// Vendor's published rating reflects the review.
	v := req(t, "GET", gw+"/api/vendors/"+vendorID, "", "", http.StatusOK)
	if v["ratingCount"].(float64) != 1 || v["ratingAvg"].(float64) != 5 {
		t.Fatalf("rating not updated after review: %+v", v)
	}
}

// TestPhotoUploadRejectsSVG locks in the MIME byte-sniffing security fix.
func TestPhotoUploadRejectsSVG(t *testing.T) {
	s := bootStack(t)
	gw := s.GatewayURL
	tok := login(t, gw, "v2@x.kz", "password123", "V2", "vendor")
	req(t, "PUT", gw+"/api/me/vendor",
		`{"name":"x","category":"c","city":"Almaty","priceFrom":1}`, tok, http.StatusOK)

	r, _ := http.NewRequest("POST", gw+"/api/me/vendor/photos",
		strings.NewReader(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`))
	r.Header.Set("Authorization", "Bearer "+tok)
	r.Header.Set("Content-Type", "image/jpeg") // forged
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("svg-with-forged-content-type must be 400, got %d", resp.StatusCode)
	}
}

// --- helpers -----------------------------------------------------------------

// login signs up + logs in a fresh user, returns the access token.
func login(t *testing.T, gw, email, pass, name, role string) string {
	t.Helper()
	req(t, "POST", gw+"/api/signup",
		`{"email":"`+email+`","password":"`+pass+`","name":"`+name+`","role":"`+role+`"}`,
		"", http.StatusCreated)
	out := req(t, "POST", gw+"/api/login",
		`{"email":"`+email+`","password":"`+pass+`"}`,
		"", http.StatusOK)
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatal("missing access token")
	}
	return tok
}

// mustLogin logs in an existing user (for admin, seeded by auth-svc).
func mustLogin(t *testing.T, gw, email, pass string) string {
	t.Helper()
	out := req(t, "POST", gw+"/api/login",
		`{"email":"`+email+`","password":"`+pass+`"}`,
		"", http.StatusOK)
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatal("missing access token")
	}
	return tok
}
