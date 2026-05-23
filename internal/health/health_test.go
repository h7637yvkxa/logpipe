package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourorg/logpipe/internal/health"
)

func TestNew_EmptySnapshot(t *testing.T) {
	c := health.New()
	if snap := c.Snapshot(); len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestSetHealthy(t *testing.T) {
	c := health.New()
	c.SetHealthy("svc-a")
	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if !snap[0].Healthy {
		t.Error("expected source to be healthy")
	}
	if snap[0].Source != "svc-a" {
		t.Errorf("unexpected source name: %s", snap[0].Source)
	}
}

func TestSetUnhealthy(t *testing.T) {
	c := health.New()
	c.SetUnhealthy("svc-b", "connection refused")
	snap := c.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].Healthy {
		t.Error("expected source to be unhealthy")
	}
	if snap[0].ErrorMsg != "connection refused" {
		t.Errorf("unexpected error msg: %s", snap[0].ErrorMsg)
	}
}

func TestHTTPHandler_AllHealthy(t *testing.T) {
	c := health.New()
	c.SetHealthy("svc-a")
	c.SetHealthy("svc-b")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c.HTTPHandler()(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ok, _ := body["ok"].(bool); !ok {
		t.Error("expected ok=true")
	}
}

func TestHTTPHandler_SomeUnhealthy(t *testing.T) {
	c := health.New()
	c.SetHealthy("svc-a")
	c.SetUnhealthy("svc-b", "timeout")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c.HTTPHandler()(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if ok, _ := body["ok"].(bool); ok {
		t.Error("expected ok=false")
	}
}
