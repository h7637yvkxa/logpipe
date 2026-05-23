package buffer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func seedBuffer(cap int, entries []Entry) *RingBuffer {
	rb := New(cap)
	for _, e := range entries {
		rb.Push(e)
	}
	return rb
}

func decodeLines(t *testing.T, data string) []map[string]string {
	t.Helper()
	var out []map[string]string
	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		if line == "" {
			continue
		}
		var m map[string]string
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("decode error: %v", err)
		}
		out = append(out, m)
	}
	return out
}

func TestReplay_AllEntries(t *testing.T) {
	entries := []Entry{
		{Source: "svcA", Line: "hello"},
		{Source: "svcB", Line: "world"},
	}
	rb := seedBuffer(10, entries)

	var buf bytes.Buffer
	if err := Replay(rb, &buf, ReplayOptions{}); err != nil {
		t.Fatal(err)
	}

	got := decodeLines(t, buf.String())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestReplay_FilterBySource(t *testing.T) {
	entries := []Entry{
		{Source: "svcA", Line: "a1"},
		{Source: "svcB", Line: "b1"},
		{Source: "svcA", Line: "a2"},
	}
	rb := seedBuffer(10, entries)

	var buf bytes.Buffer
	if err := Replay(rb, &buf, ReplayOptions{Source: "svcA"}); err != nil {
		t.Fatal(err)
	}

	got := decodeLines(t, buf.String())
	if len(got) != 2 {
		t.Fatalf("expected 2 svcA entries, got %d", len(got))
	}
	for _, m := range got {
		if m["source"] != "svcA" {
			t.Errorf("unexpected source %q", m["source"])
		}
	}
}

func TestReplay_Limit(t *testing.T) {
	entries := []Entry{
		{Source: "s", Line: "l1"},
		{Source: "s", Line: "l2"},
		{Source: "s", Line: "l3"},
		{Source: "s", Line: "l4"},
	}
	rb := seedBuffer(10, entries)

	var buf bytes.Buffer
	if err := Replay(rb, &buf, ReplayOptions{Limit: 2}); err != nil {
		t.Fatal(err)
	}

	got := decodeLines(t, buf.String())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0]["line"] != "l3" || got[1]["line"] != "l4" {
		t.Errorf("expected last 2 lines, got %v", got)
	}
}

func TestReplay_EmptyBuffer(t *testing.T) {
	rb := New(5)
	var buf bytes.Buffer
	if err := Replay(rb, &buf, ReplayOptions{}); err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output, got %q", buf.String())
	}
}
