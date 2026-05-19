package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logpipe/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logpipe-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `
sources:
  - name: app
    path: /var/log/app.log
    format: json
output:
  type: stdout
filter:
  levels: [error, warn]
`
	cfg, err := config.Load(writeTemp(t, raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Sources) != 1 {
		t.Errorf("expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "app" {
		t.Errorf("expected source name 'app', got %q", cfg.Sources[0].Name)
	}
	if cfg.Output.Type != "stdout" {
		t.Errorf("expected output type 'stdout', got %q", cfg.Output.Type)
	}
}

func TestLoad_MissingSources(t *testing.T) {
	raw := `output:\n  type: stdout\n`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing sources, got nil")
	}
}

func TestLoad_InvalidFormat(t *testing.T) {
	raw := `
sources:
  - name: svc
    path: /tmp/svc.log
    format: xml
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestLoad_FileOutputMissingPath(t *testing.T) {
	raw := `
sources:
  - name: svc
    path: /tmp/svc.log
output:
  type: file
`
	_, err := config.Load(writeTemp(t, raw))
	if err == nil {
		t.Fatal("expected error for missing output path, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
