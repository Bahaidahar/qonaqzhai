package chat_test

import (
	"context"
	"errors"
	"testing"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/chat"
	"qonaqzhai-backend/internal/usecase/inmem"
)

type stubAI struct {
	called int
	reply  *usecase.ChatReply
	err    error
}

func (s *stubAI) Generate(_ context.Context, _ string, _ []usecase.VendorRef) (*usecase.ChatReply, error) {
	s.called++
	return s.reply, s.err
}

func TestGenerateRejectsEmptyMessage(t *testing.T) {
	t.Parallel()
	svc := chat.New(chat.Deps{Vendors: inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)})
	if _, err := svc.Generate(context.Background(), "  "); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err=%v want ErrInvalidInput", err)
	}
}

func TestGenerateUsesAIWhenAvailable(t *testing.T) {
	t.Parallel()
	ai := &stubAI{reply: &usecase.ChatReply{Reply: "hi from gemini"}}
	svc := chat.New(chat.Deps{
		Vendors: inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New),
		AI:      ai,
	})
	r, err := svc.Generate(context.Background(), "plan my toi")
	if err != nil {
		t.Fatal(err)
	}
	if r.Reply != "hi from gemini" {
		t.Errorf("reply=%q", r.Reply)
	}
	if ai.called != 1 {
		t.Errorf("ai called %d times", ai.called)
	}
}

func TestGenerateFallsBackOnAIError(t *testing.T) {
	t.Parallel()
	ai := &stubAI{err: errors.New("network")}
	svc := chat.New(chat.Deps{
		Vendors: inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New),
		AI:      ai,
	})
	r, err := svc.Generate(context.Background(), "budget please")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Blocks) == 0 || r.Blocks[0].Type != "budget" {
		t.Errorf("expected budget fallback, got %+v", r.Blocks)
	}
}

func TestFallbackKeywords(t *testing.T) {
	t.Parallel()
	svc := chat.New(chat.Deps{Vendors: inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)})
	cases := []struct {
		msg, wantType string
	}{
		{"бюджет на свадьбу", "budget"},
		{"найди фотографа", "vendors"},
		{"hello", "plan"},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			r, err := svc.Generate(context.Background(), c.msg)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Blocks) == 0 || r.Blocks[0].Type != c.wantType {
				t.Errorf("got %+v want type %q", r.Blocks, c.wantType)
			}
		})
	}
}
