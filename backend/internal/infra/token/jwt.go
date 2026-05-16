// Package token implements the JWT-based TokenIssuer port.
package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// JWT implements usecase.TokenIssuer using HS256-signed access tokens.
type JWT struct {
	secret []byte
}

// New returns a JWT issuer for the given secret.
func New(secret []byte) *JWT { return &JWT{secret: secret} }

type claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Issue signs a fresh access token for u with the given TTL.
func (j *JWT) Issue(u *domain.User, ttl time.Duration) (string, error) {
	c := &claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "qonaqzhai",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signed, err := tok.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// Parse verifies and decodes a JWT issued by this signer.
func (j *JWT) Parse(raw string) (usecase.Claims, error) {
	tok, err := jwt.ParseWithClaims(raw, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})
	if err != nil {
		return usecase.Claims{}, fmt.Errorf("parse jwt: %w", err)
	}
	c, ok := tok.Claims.(*claims)
	if !ok || !tok.Valid {
		return usecase.Claims{}, errors.New("invalid token")
	}
	return usecase.Claims{
		UserID: c.UserID,
		Email:  c.Email,
		Role:   domain.Role(c.Role),
	}, nil
}
