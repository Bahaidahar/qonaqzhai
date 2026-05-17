// Binary realtime-svc owns the WebSocket fan-out hub.
// JWT auth → auth-svc gRPC. Message persistence → core-svc REST /internal endpoint.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "qonaqzhai-backend/gen/proto/auth/v1"
	"qonaqzhai-backend/internal/adapter/http/handler"
	"qonaqzhai-backend/internal/adapter/http/httpx"
	wshub "qonaqzhai-backend/internal/adapter/ws"
	"qonaqzhai-backend/internal/infra/config"
	"qonaqzhai-backend/internal/infra/logger"
	"qonaqzhai-backend/internal/infra/token"
)

func main() {
	_ = godotenv.Load()
	if _, err := config.Load(); err != nil {
		log.Fatalf("config: %v", err)
	}

	httpAddr := envOr("REALTIME_HTTP_ADDR", ":8083")
	authGRPC := envOr("AUTH_GRPC_TARGET", "localhost:9091")
	coreBase := envOr("CORE_INTERNAL_URL", "http://localhost:8082")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logr := logger.New("info")

	authConn, err := grpc.NewClient(authGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial auth: %v", err)
	}
	defer authConn.Close()
	tokVerifier := token.NewRemote(ctx, authv1.NewAuthServiceClient(authConn))

	hub := wshub.NewHub(logr)
	core := &coreClient{base: coreBase, http: &http.Client{Timeout: 5 * time.Second}}

	// On incoming WS message: persist via core, then fan out the persisted form.
	hub.OnIncoming = func(senderID string, m wshub.IncomingMessage) {
		out, err := core.persistThreadMessage(ctx, m.ThreadID, senderID, m.Text)
		if err != nil {
			logr.Warn("persist thread msg", "err", err.Error(), "thread", m.ThreadID)
			return
		}
		hub.SendToUsers(wshub.Envelope{Op: "thread.message", Data: out.Message}, out.Participants...)
	}

	wsHandler := handler.NewWS(hub, tokVerifier)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "realtime"})
	})
	mux.HandleFunc("GET /api/ws", wsHandler.Connect)

	srv := &http.Server{
		Addr:              httpAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		// no WriteTimeout — websockets are long-lived.
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		log.Printf("realtime-svc listening on %s · auth=%s · core=%s", httpAddr, authGRPC, coreBase)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("realtime-svc shutting down")
	shutCtx, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()
	_ = srv.Shutdown(shutCtx)
}

type coreClient struct {
	base string
	http *http.Client
}

type threadMessage struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"threadId"`
	SenderID  string    `json:"senderId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}

type persistResult struct {
	Message      *threadMessage `json:"message"`
	Participants []string       `json:"participants"`
}

func (c *coreClient) persistThreadMessage(ctx context.Context, threadID, senderID, text string) (*persistResult, error) {
	body, _ := json.Marshal(map[string]string{"senderId": senderID, "text": text})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/internal/threads/%s/messages", c.base, threadID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		raw, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("core %d: %s", res.StatusCode, string(raw))
	}
	var out persistResult
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
