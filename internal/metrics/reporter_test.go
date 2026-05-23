package metrics

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"
)

func TestReporter_WritesSnapshot(t *testing.T) {
	c := New()
	c.IncRead("svc")
	c.IncWritten()

	var buf bytes.Buffer
	r := NewReporter(c, 50*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()
	r.Run(ctx)

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) == 0 {
		t.Fatal("expected at least one metrics line")
	}

	var snap map[string]int64
	if err := json.Unmarshal(lines[0], &snap); err != nil {
		t.Fatalf("invalid JSON in metrics output: %v", err)
	}
	if snap["lines_read"] != 1 {
		t.Errorf("lines_read: want 1, got %d", snap["lines_read"])
	}
	if snap["lines_written"] != 1 {
		t.Errorf("lines_written: want 1, got %d", snap["lines_written"])
	}
	if _, ok := snap["_ts"]; !ok {
		t.Error("snapshot missing _ts field")
	}
}

func TestReporter_FinalFlushOnCancel(t *testing.T) {
	c := New()
	c.IncParseError()

	var buf bytes.Buffer
	r := NewReporter(c, 10*time.Second, &buf) // long interval — only cancel flush fires

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	r.Run(ctx)

	if buf.Len() == 0 {
		t.Fatal("expected final flush on context cancellation")
	}

	var snap map[string]int64
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &snap); err != nil {
		t.Fatalf("invalid JSON on final flush: %v", err)
	}
	if snap["parse_errors"] != 1 {
		t.Errorf("parse_errors: want 1, got %d", snap["parse_errors"])
	}
}

func TestNewReporter_DefaultsToStderr(t *testing.T) {
	c := New()
	r := NewReporter(c, time.Second, nil)
	if r.out == nil {
		t.Error("expected non-nil writer when nil passed to NewReporter")
	}
}
