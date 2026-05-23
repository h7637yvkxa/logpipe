package dedupe

import (
	"testing"
	"time"
)

func collect(ch <-chan Entry) []Entry {
	var results []Entry
	for e := range ch {
		results = append(results, e)
	}
	return results
}

func TestPipeline_UniqueEntriesPassThrough(t *testing.T) {
	p := NewPipeline(5 * time.Second)
	in := make(chan Entry, 3)
	in <- Entry{"msg": "a"}
	in <- Entry{"msg": "b"}
	in <- Entry{"msg": "c"}
	close(in)

	results := collect(p.Run(in))
	if len(results) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(results))
	}
}

func TestPipeline_DuplicatesFiltered(t *testing.T) {
	p := NewPipeline(5 * time.Second)
	in := make(chan Entry, 4)
	in <- Entry{"msg": "hello"}
	in <- Entry{"msg": "hello"}
	in <- Entry{"msg": "world"}
	in <- Entry{"msg": "hello"}
	close(in)

	results := collect(p.Run(in))
	if len(results) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(results))
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	p := NewPipeline(5 * time.Second)
	in := make(chan Entry)
	close(in)

	results := collect(p.Run(in))
	if len(results) != 0 {
		t.Fatalf("expected 0 entries from empty input, got %d", len(results))
	}
}

func TestPipeline_OutputChannelClosedAfterInput(t *testing.T) {
	p := NewPipeline(time.Second)
	in := make(chan Entry, 1)
	in <- Entry{"x": "1"}
	close(in)

	out := p.Run(in)
	collect(out)
	// Verify channel is closed (collect drains it; no panic means success)
}
