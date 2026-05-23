package sampling

import (
	"testing"

	"github.com/your/logpipe/internal/normalize"
)

func feedPipeline(entries []normalize.Entry, s *Sampler) []normalize.Entry {
	in := make(chan normalize.Entry, len(entries))
	for _, e := range entries {
		in <- e
	}
	close(in)
	out := Pipeline(in, s)
	var result []normalize.Entry
	for e := range out {
		result = append(result, e)
	}
	return result
}

func TestPipeline_AllKept(t *testing.T) {
	s := New(Config{Rate: 1.0})
	input := []normalize.Entry{
		{Source: "a", Message: "one"},
		{Source: "a", Message: "two"},
		{Source: "a", Message: "three"},
	}
	got := feedPipeline(input, s)
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestPipeline_AllDropped(t *testing.T) {
	s := New(Config{Rate: 0.0})
	input := []normalize.Entry{
		{Source: "b", Message: "x"},
		{Source: "b", Message: "y"},
	}
	got := feedPipeline(input, s)
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	s := New(Config{Rate: 1.0})
	got := feedPipeline(nil, s)
	if len(got) != 0 {
		t.Fatalf("expected empty output, got %d entries", len(got))
	}
}

func TestPipeline_OutputClosedAfterInput(t *testing.T) {
	s := New(Config{Rate: 1.0})
	in := make(chan normalize.Entry)
	close(in)
	out := Pipeline(in, s)
	_, open := <-out
	if open {
		t.Fatal("expected output channel to be closed")
	}
}
