package health_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/health"
)

func freePort() int {
	// Use a fixed high port range for tests; real projects use net.Listen trick.
	return 19876
}

func TestServer_StartsAndResponds(t *testing.T) {
	port := freePort()
	checker := health.New()
	checker.SetHealthy("test-svc")

	cfg := health.ServerConfig{Port: port}
	srv := health.NewServer(cfg, checker)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		t.Fatalf("server failed to start: %v", err)
	}

	// Give the server a moment to bind.
	time.Sleep(80 * time.Millisecond)

	url := fmt.Sprintf("http://localhost:%d/health", port)
	resp, err := http.Get(url) //nolint:noctx
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestServer_ShutdownOnContextCancel(t *testing.T) {
	port := freePort() + 1
	checker := health.New()

	cfg := health.ServerConfig{Port: port}
	srv := health.NewServer(cfg, checker)

	ctx, cancel := context.WithCancel(context.Background())
	if err := srv.Start(ctx); err != nil {
		t.Fatalf("server failed to start: %v", err)
	}
	time.Sleep(80 * time.Millisecond)

	cancel()
	time.Sleep(200 * time.Millisecond)

	// After shutdown, connections should be refused.
	_, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port)) //nolint:noctx
	if err == nil {
		t.Error("expected connection error after shutdown, got nil")
	}
}
