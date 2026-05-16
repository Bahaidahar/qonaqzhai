// Package push implements the Pusher port against Firebase Cloud Messaging HTTP v1.
package push

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"qonaqzhai-backend/internal/usecase"
)

// Config bundles FCM credentials.
type Config struct {
	ProjectID         string
	ServiceAccountKey []byte // contents of the JSON key file
}

// FCM implements usecase.Pusher.
type FCM struct {
	projectID string
	signer    *serviceAccountSigner
	client    *http.Client

	mu   sync.Mutex
	tok  string
	exp  time.Time
}

// New returns nil when configuration is incomplete (push disabled).
func New(cfg Config) (*FCM, error) {
	if cfg.ProjectID == "" || len(cfg.ServiceAccountKey) == 0 {
		return nil, nil
	}
	signer, err := parseServiceAccount(cfg.ServiceAccountKey)
	if err != nil {
		return nil, err
	}
	return &FCM{
		projectID: cfg.ProjectID,
		signer:    signer,
		client:    &http.Client{Timeout: 15 * time.Second},
	}, nil
}

// Send delivers a notification to each token.
// Failures on individual tokens are logged but do not abort the loop.
func (f *FCM) Send(ctx context.Context, tokens []string, title, body string, data map[string]string) error {
	if f == nil {
		return errors.New("fcm disabled")
	}
	if len(tokens) == 0 {
		return nil
	}
	tok, err := f.accessToken(ctx)
	if err != nil {
		return err
	}
	endpoint := fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", f.projectID)

	for _, target := range tokens {
		payload := map[string]any{
			"message": map[string]any{
				"token": target,
				"notification": map[string]string{
					"title": title,
					"body":  body,
				},
				"data": stringMap(data),
			},
		}
		buf, _ := json.Marshal(payload)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(buf))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+tok)
		req.Header.Set("Content-Type", "application/json")
		res, err := f.client.Do(req)
		if err != nil {
			return fmt.Errorf("fcm send: %w", err)
		}
		if res.StatusCode >= 300 {
			b, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
			res.Body.Close()
			return fmt.Errorf("fcm status %d: %s", res.StatusCode, strings.TrimSpace(string(b)))
		}
		res.Body.Close()
	}
	return nil
}

func stringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func (f *FCM) accessToken(ctx context.Context) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.tok != "" && time.Now().Before(f.exp.Add(-time.Minute)) {
		return f.tok, nil
	}
	jwtTok, err := f.signer.sign(time.Now())
	if err != nil {
		return "", err
	}
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", jwtTok)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token exchange: %w", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		return "", fmt.Errorf("token exchange status %d: %s", res.StatusCode, string(body))
	}
	var tr struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", err
	}
	f.tok = tr.AccessToken
	f.exp = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return f.tok, nil
}

type serviceAccountSigner struct {
	clientEmail string
	privateKey  *rsa.PrivateKey
}

func parseServiceAccount(raw []byte) (*serviceAccountSigner, error) {
	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	if err := json.Unmarshal(raw, &sa); err != nil {
		return nil, fmt.Errorf("parse service account: %w", err)
	}
	if sa.ClientEmail == "" || sa.PrivateKey == "" {
		return nil, errors.New("service account JSON missing client_email or private_key")
	}
	block, _ := pem.Decode([]byte(sa.PrivateKey))
	if block == nil {
		return nil, errors.New("invalid PEM private_key")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse pkcs8: %w", err)
	}
	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("service account private_key is not RSA")
	}
	return &serviceAccountSigner{clientEmail: sa.ClientEmail, privateKey: rsaKey}, nil
}

func (s *serviceAccountSigner) sign(now time.Time) (string, error) {
	claims := jwt.MapClaims{
		"iss":   s.clientEmail,
		"scope": "https://www.googleapis.com/auth/firebase.messaging",
		"aud":   "https://oauth2.googleapis.com/token",
		"iat":   now.Unix(),
		"exp":   now.Add(time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return tok.SignedString(s.privateKey)
}

var _ usecase.Pusher = (*FCM)(nil)
