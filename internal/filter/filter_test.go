package filter

import (
	"testing"
)

func entry(kv ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestMatch_NoRules(t *testing.T) {
	f := New(nil)
	if !f.Match(entry("level", "error")) {
		t.Fatal("empty rule set should match everything")
	}
}

func TestMatch_FieldContains(t *testing.T) {
	f := New([]Rule{{Field: "level", Match: "error"}})
	if !f.Match(entry("level", "error")) {
		t.Fatal("expected match")
	}
	if f.Match(entry("level", "info")) {
		t.Fatal("expected no match")
	}
}

func TestMatch_CaseInsensitive(t *testing.T) {
	f := New([]Rule{{Field: "level", Match: "ERROR"}})
	if !f.Match(entry("level", "error")) {
		t.Fatal("match should be case-insensitive")
	}
}

func TestMatch_Invert(t *testing.T) {
	f := New([]Rule{{Field: "level", Match: "debug", Invert: true}})
	if !f.Match(entry("level", "info")) {
		t.Fatal("inverted rule: non-matching value should pass")
	}
	if f.Match(entry("level", "debug")) {
		t.Fatal("inverted rule: matching value should be excluded")
	}
}

func TestMatch_MissingField(t *testing.T) {
	f := New([]Rule{{Field: "service", Match: "api"}})
	if f.Match(entry("level", "info")) {
		t.Fatal("missing required field should not match")
	}
}

func TestMatch_MissingField_Inverted(t *testing.T) {
	f := New([]Rule{{Field: "service", Match: "api", Invert: true}})
	// field absent + inverted → rule skipped → passes
	if !f.Match(entry("level", "info")) {
		t.Fatal("missing field with inverted rule should pass")
	}
}

func TestPipeline_FiltersEntries(t *testing.T) {
	f := New([]Rule{{Field: "level", Match: "error"}})

	in := make(chan map[string]interface{}, 4)
	in <- entry("level", "error", "msg", "boom")
	in <- entry("level", "info", "msg", "ok")
	in <- entry("level", "error", "msg", "fail")
	close(in)

	out := Pipeline(in, f)
	var got []map[string]interface{}
	for e := range out {
		got = append(got, e)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}
