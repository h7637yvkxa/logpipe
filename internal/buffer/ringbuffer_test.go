package buffer

import (
	"fmt"
	"testing"
)

func TestNew_DefaultsToOneWhenZero(t *testing.T) {
	rb := New(0)
	if rb.Cap() != 1 {
		t.Fatalf("expected cap 1, got %d", rb.Cap())
	}
}

func TestPush_BelowCapacity(t *testing.T) {
	rb := New(5)
	rb.Push(Entry{Source: "svc", Line: "line1"})
	rb.Push(Entry{Source: "svc", Line: "line2"})

	if rb.Len() != 2 {
		t.Fatalf("expected len 2, got %d", rb.Len())
	}
}

func TestSnapshot_OrderPreserved(t *testing.T) {
	rb := New(4)
	for i := 0; i < 4; i++ {
		rb.Push(Entry{Source: "s", Line: fmt.Sprintf("line%d", i)})
	}

	snap := rb.Snapshot()
	for i, e := range snap {
		want := fmt.Sprintf("line%d", i)
		if e.Line != want {
			t.Errorf("pos %d: got %q, want %q", i, e.Line, want)
		}
	}
}

func TestPush_OverwritesOldest(t *testing.T) {
	rb := New(3)
	for i := 0; i < 5; i++ {
		rb.Push(Entry{Source: "s", Line: fmt.Sprintf("line%d", i)})
	}

	if rb.Len() != 3 {
		t.Fatalf("expected len 3, got %d", rb.Len())
	}

	snap := rb.Snapshot()
	expected := []string{"line2", "line3", "line4"}
	for i, e := range snap {
		if e.Line != expected[i] {
			t.Errorf("pos %d: got %q, want %q", i, e.Line, expected[i])
		}
	}
}

func TestSnapshot_EmptyBuffer(t *testing.T) {
	rb := New(10)
	snap := rb.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestPush_Concurrent(t *testing.T) {
	rb := New(50)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 20; j++ {
				rb.Push(Entry{Source: "s", Line: fmt.Sprintf("%d-%d", n, j)})
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if rb.Len() > 50 {
		t.Fatalf("len %d exceeds cap 50", rb.Len())
	}
}
