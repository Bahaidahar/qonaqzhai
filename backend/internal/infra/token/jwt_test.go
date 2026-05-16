package token_test

import (
	"strings"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/infra/token"
)

func TestIssueAndParseRoundtrip(t *testing.T) {
	t.Parallel()
	j := token.New([]byte("test-secret"))
	u := &domain.User{ID: "u-1", Email: "x@y.com", Role: domain.RoleCustomer}
	tok, err := j.Issue(u, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(tok, ".") {
		t.Fatalf("token doesn't look like JWT: %q", tok)
	}
	c, err := j.Parse(tok)
	if err != nil {
		t.Fatal(err)
	}
	if c.UserID != "u-1" || c.Email != "x@y.com" || c.Role != domain.RoleCustomer {
		t.Errorf("claims=%+v", c)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	t.Parallel()
	tok, _ := token.New([]byte("a")).Issue(&domain.User{ID: "u", Email: "e", Role: domain.RoleCustomer}, time.Hour)
	if _, err := token.New([]byte("b")).Parse(tok); err == nil {
		t.Error("token verified with wrong secret")
	}
}

func TestParseRejectsExpired(t *testing.T) {
	t.Parallel()
	j := token.New([]byte("s"))
	tok, _ := j.Issue(&domain.User{ID: "u", Email: "e", Role: domain.RoleCustomer}, -time.Minute)
	if _, err := j.Parse(tok); err == nil {
		t.Error("expired token verified")
	}
}
