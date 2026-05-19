package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/tail"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "logpipe-tail-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	path := writeTempLog(t, "")

	out := make(chan tail.Line, 10)
	tr := tail.New(path, "svc-a", out)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		_ = tr.Run(ctx)
	}()

	// Give the tailer time to open and seek.
	time.Sleep(50 * time.Millisecond)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open log file for writing: %v", err)
	}
	_, _ = f.WriteString(`{"level":"info","msg":"hello"}` + "\n")
	f.Close()

	select {
	case line := <-out:
		if line.Source != "svc-a" {
			t.Errorf("expected source svc-a, got %s", line.Source)
		}
		if line.Text == "" {
			t.Error("expected non-empty line text")
		}
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for tailed line")
	}

	cancel()
}

func TestTailer_MissingFile(t *testing.T) {
	out := make(chan tail.Line, 1)
	tr := tail.New("/nonexistent/path/file.log", "svc-b", out)
	err := tr.Run(context.Background())
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
