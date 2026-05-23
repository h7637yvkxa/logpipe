package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_NoLimit(t *testing.T) {
	l := New(0)
	for i := 0; i < 1000; i++ {
		if !l.Allow("svc") {
			t.Fatal("expected all entries allowed when rate is 0")
		}
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	l := New(5)
	for i := 0; i < 5; i++ {
		if !l.Allow("svc") {
			t.Fatalf("entry %d should be allowed within limit", i)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(3)
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow("svc") {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed, got %d", allowed)
	}
}

func TestAllow_IndependentSources(t *testing.T) {
	l := New(2)
	if !l.Allow("a") { t.Fatal("a[0] should be allowed") }
	if !l.Allow("b") { t.Fatal("b[0] should be allowed") }
	if !l.Allow("a") { t.Fatal("a[1] should be allowed") }
	if !l.Allow("b") { t.Fatal("b[1] should be allowed") }
	if l.Allow("a")  { t.Fatal("a[2] should be denied") }
	if l.Allow("b")  { t.Fatal("b[2] should be denied") }
}

func TestAllow_WindowResets(t *testing.T) {
	l := New(2)
	l.Allow("svc")
	l.Allow("svc")
	if l.Allow("svc") {
		t.Fatal("third entry should be denied before window resets")
	}

	// Manually rewind window start to simulate time passing.
	l.mu.Lock()
	l.buckets["svc"].windowStart = time.Now().Add(-2 * time.Second)
	l.mu.Unlock()

	if !l.Allow("svc") {
		t.Fatal("first entry after window reset should be allowed")
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := New(1)
	l.Allow("svc")
	if l.Allow("svc") {
		t.Fatal("second entry should be denied")
	}
	l.Reset()
	if !l.Allow("svc") {
		t.Fatal("entry after Reset should be allowed")
	}
}
