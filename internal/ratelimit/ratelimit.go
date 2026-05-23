package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of log entries per second per source.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	buckets  map[string]*bucket
}

type bucket struct {
	count    int
	windowStart time.Time
}

// New creates a Limiter that allows up to rate entries per second per source.
// A rate of 0 disables limiting.
func New(rate int) *Limiter {
	return &Limiter{
		rate:    rate,
		buckets: make(map[string]*bucket),
	}
}

// Allow returns true if the entry from the given source should be passed
// through, or false if it exceeds the configured rate limit.
func (l *Limiter) Allow(source string) bool {
	if l.rate <= 0 {
		return true
	}

	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[source]
	if !ok {
		l.buckets[source] = &bucket{count: 1, windowStart: now}
		return true
	}

	if now.Sub(b.windowStart) >= time.Second {
		b.count = 1
		b.windowStart = now
		return true
	}

	if b.count >= l.rate {
		return false
	}

	b.count++
	return true
}

// Reset clears the rate limit state for all sources.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string]*bucket)
}
