package output_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
	"github.com/yourorg/logpipe/internal/output"
)

func TestWriter_StdoutType(t *testing.T) {
	w, err := output.New(config.OutputConfig{Type: "stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriter_FileType_MissingPath(t *testing.T) {
	_, err := output.New(config.OutputConfig{Type: "file", Path: ""})
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestWriter_FileType_WritesJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	w, err := output.New(config.OutputConfig{Type: "file", Path: path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := map[string]interface{}{
		"level":   "info",
		"message": "hello",
		"service": "svc-a",
	}
	if err := w.Write(entry); err != nil {
		t.Fatalf("write error: %v", err)
	}
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal(bytes.TrimSpace(data), &got); err != nil {
		t.Fatalf("invalid JSON in output: %v", err)
	}
	if got["message"] != "hello" {
		t.Errorf("expected message 'hello', got %v", got["message"])
	}
	if got["level"] != "info" {
		t.Errorf("expected level 'info', got %v", got["level"])
	}
}

func TestWriter_Write_InvalidEntry(t *testing.T) {
	w, err := output.New(config.OutputConfig{Type: "stdout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()

	// json.Marshal cannot handle channel values
	bad := map[string]interface{}{
		"ch": make(chan int),
	}
	if err := w.Write(bad); err == nil {
		t.Error("expected marshalling error for invalid entry, got nil")
	}
}
