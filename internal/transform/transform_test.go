package transform

import (
	"testing"

	"github.com/user/logpipe/internal/normalize"
)

func baseEntry() normalize.Entry {
	return normalize.Entry{
		Timestamp: "2024-01-01T00:00:00Z",
		Level:     "info",
		Message:   "Hello World",
		Source:    "svc",
		Extra:     map[string]any{"env": "production", "code": 200},
	}
}

func TestApply_NoRules(t *testing.T) {
	tr := New(nil)
	e := baseEntry()
	out := tr.Apply(e)
	if out.Level != e.Level || out.Message != e.Message {
		t.Fatalf("expected unchanged entry, got %+v", out)
	}
}

func TestApply_Uppercase(t *testing.T) {
	tr := New([]Rule{{Field: "level", Op: "uppercase"}})
	out := tr.Apply(baseEntry())
	if out.Level != "INFO" {
		t.Fatalf("expected INFO, got %q", out.Level)
	}
}

func TestApply_Lowercase(t *testing.T) {
	tr := New([]Rule{{Field: "message", Op: "lowercase"}})
	out := tr.Apply(baseEntry())
	if out.Message != "hello world" {
		t.Fatalf("expected 'hello world', got %q", out.Message)
	}
}

func TestApply_Truncate(t *testing.T) {
	tr := New([]Rule{{Field: "message", Op: "truncate", Arg: "5"}})
	out := tr.Apply(baseEntry())
	if out.Message != "Hello" {
		t.Fatalf("expected 'Hello', got %q", out.Message)
	}
}

func TestApply_TruncateBelowLength(t *testing.T) {
	tr := New([]Rule{{Field: "message", Op: "truncate", Arg: "100"}})
	out := tr.Apply(baseEntry())
	if out.Message != "Hello World" {
		t.Fatalf("expected unchanged message, got %q", out.Message)
	}
}

func TestApply_Replace(t *testing.T) {
	tr := New([]Rule{{Field: "message", Op: "replace", Arg: "World:Go"}})
	out := tr.Apply(baseEntry())
	if out.Message != "Hello Go" {
		t.Fatalf("expected 'Hello Go', got %q", out.Message)
	}
}

func TestApply_ExtraStringField(t *testing.T) {
	tr := New([]Rule{{Field: "env", Op: "uppercase"}})
	out := tr.Apply(baseEntry())
	if out.Extra["env"] != "PRODUCTION" {
		t.Fatalf("expected PRODUCTION, got %v", out.Extra["env"])
	}
}

func TestApply_ExtraNonStringFieldUnchanged(t *testing.T) {
	tr := New([]Rule{{Field: "code", Op: "uppercase"}})
	out := tr.Apply(baseEntry())
	if out.Extra["code"] != 200 {
		t.Fatalf("expected 200 unchanged, got %v", out.Extra["code"])
	}
}

func TestPipeline_TransformsEntries(t *testing.T) {
	in := make(chan normalize.Entry, 2)
	in <- baseEntry()
	in <- baseEntry()
	close(in)

	tr := New([]Rule{{Field: "level", Op: "uppercase"}})
	out := Pipeline(in, tr)

	count := 0
	for e := range out {
		if e.Level != "INFO" {
			t.Fatalf("expected INFO, got %q", e.Level)
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 entries, got %d", count)
	}
}

func TestPipeline_ClosedOnInputClose(t *testing.T) {
	in := make(chan normalize.Entry)
	close(in)
	tr := New(nil)
	out := Pipeline(in, tr)
	_, open := <-out
	if open {
		t.Fatal("expected output channel to be closed")
	}
}
