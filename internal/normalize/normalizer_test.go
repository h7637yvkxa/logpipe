package normalize

import (
	"testing"
	"time"
)

func TestNormalize_StandardFields(t *testing.T) {
	n := New("svc-a")
	line := `{"timestamp":"2024-01-15T10:00:00Z","level":"error","message":"something broke","code":42}`

	entry, err := n.Normalize(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Service != "svc-a" {
		t.Errorf("expected service svc-a, got %s", entry.Service)
	}
	if entry.Level != "error" {
		t.Errorf("expected level error, got %s", entry.Level)
	}
	if entry.Message != "something broke" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
	if entry.Fields["code"] == nil {
		t.Error("expected extra field 'code' to be present")
	}
}

func TestNormalize_AliasFields(t *testing.T) {
	n := New("svc-b")
	line := `{"ts":1705312800,"lvl":"warn","msg":"disk full"}`

	entry, err := n.Normalize(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "warn" {
		t.Errorf("expected warn, got %s", entry.Level)
	}
	if entry.Message != "disk full" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
	expected := time.Unix(1705312800, 0).UTC()
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entry.Timestamp)
	}
}

func TestNormalize_DefaultsWhenMissing(t *testing.T) {
	n := New("svc-c")
	line := `{"msg":"hello world"}`

	entry, err := n.Normalize(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "info" {
		t.Errorf("expected default level info, got %s", entry.Level)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNormalize_InvalidJSON(t *testing.T) {
	n := New("svc-d")
	_, err := n.Normalize(`not json at all`)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestNormalize_NoExtraFields(t *testing.T) {
	n := New("svc-e")
	line := `{"level":"info","message":"clean"}`

	entry, err := n.Normalize(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Fields != nil {
		t.Errorf("expected nil fields, got %v", entry.Fields)
	}
}
