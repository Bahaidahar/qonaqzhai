package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSigner issues and parses HS256 JWTs. Only auth-svc should construct one
// (it owns the secret).
type JWTSigner struct {
	secret []byte
	issuer string
}

// NewJWTSigner returns a signer with the given HMAC secret.
func NewJWTSigner(secret []byte, issuer string) *JWTSigner {
	if issuer == "" {
		issuer = "qonaqzhai"
	}
	return &JWTSigner{secret: secret, issuer: issuer}
}

type accessClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Status string `json:"status,omitempty"`
	jwt.RegisteredClaims
}

// Issue signs a fresh access token for the given principal.
func (s *JWTSigner) Issue(c Claims, ttl time.Duration) (string, error) {
	now := time.Now()
	cl := &accessClaims{
		UserID: c.UserID,
		Email:  c.Email,
		Role:   c.Role,
		Status: c.Status,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Issuer:    s.issuer,
			Subject:   c.UserID,
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// Parse verifies a token issued by this signer.
func (s *JWTSigner) Parse(raw string) (Claims, time.Time, error) {
	tok, err := jwt.ParseWithClaims(raw, &accessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return Claims{}, time.Time{}, fmt.Errorf("parse jwt: %w", err)
	}
	cl, ok := tok.Claims.(*accessClaims)
	if !ok || !tok.Valid {
		return Claims{}, time.Time{}, errors.New("invalid token")
	}
	c := Claims{
		UserID: cl.UserID,
		Email:  cl.Email,
		Role:   cl.Role,
		Status: cl.Status,
	}
	exp := time.Time{}
	if cl.ExpiresAt != nil {
		exp = cl.ExpiresAt.Time
	}
	return c, exp, nil
}
