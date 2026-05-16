package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

type chatBlockResp struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

type chatRespJSON struct {
	Reply  string          `json:"reply"`
	Blocks []chatBlockResp `json:"blocks"`
}

func TestChatRequiresAuth(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("POST", "/api/chat", "", map[string]any{"message": "hi"})
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401 got %d", res.StatusCode)
	}
}

func TestChatBudgetReturnsBudgetBlock(t *testing.T) {
	e := newEnv(t)
	tok := e.signup("a@b.kz", "password123", "A", "customer")
	res, body := e.do("POST", "/api/chat", tok, map[string]any{"message": "show me a budget"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("chat: %d %s", res.StatusCode, body)
	}
	r := decode[chatRespJSON](t, body)
	if len(r.Blocks) != 1 || r.Blocks[0].Type != "budget" {
		t.Errorf("expected budget block, got %+v", r.Blocks)
	}
}

func TestChatVendorsUsesRealVendors(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Studio Aitu", "category": "Photo & Video", "city": "Almaty", "priceFrom": 450000,
	})
	v := decode[domain.Vendor](t, body)
	_, _ = e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{"status": "approved"})

	tok := e.signup("c@b.kz", "password123", "C", "customer")
	_, body = e.do("POST", "/api/chat", tok, map[string]any{"message": "find me a vendor photographer"})
	r := decode[chatRespJSON](t, body)
	if len(r.Blocks) != 1 || r.Blocks[0].Type != "vendors" {
		t.Fatalf("expected vendors block, got %+v", r.Blocks)
	}
	items, _ := r.Blocks[0].Data["items"].([]any)
	if len(items) != 1 {
		t.Errorf("expected 1 vendor in chat reply, got %d", len(items))
	}
}

func TestChatEmptyMessage(t *testing.T) {
	e := newEnv(t)
	tok := e.signup("a@b.kz", "password123", "A", "customer")
	res, _ := e.do("POST", "/api/chat", tok, map[string]any{"message": "   "})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}
