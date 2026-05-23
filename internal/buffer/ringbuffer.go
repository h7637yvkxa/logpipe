package buffer

import "sync"

// Entry holds a single log line with its source label.
type Entry struct {
	Source string
	Line   string
}

// RingBuffer is a fixed-capacity circular buffer of log entries.
// When full, the oldest entry is overwritten.
type RingBuffer struct {
	mu       sync.Mutex
	items    []Entry
	cap      int
	head     int
	count    int
}

// New creates a RingBuffer with the given capacity.
func New(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 1
	}
	return &RingBuffer{
		items: make([]Entry, capacity),
		cap:   capacity,
	}
}

// Push adds an entry to the buffer, evicting the oldest if full.
func (r *RingBuffer) Push(e Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pos := (r.head + r.count) % r.cap
	r.items[pos] = e

	if r.count < r.cap {
		r.count++
	} else {
		// overwrite: advance head
		r.head = (r.head + 1) % r.cap
	}
}

// Snapshot returns a copy of all entries in insertion order.
func (r *RingBuffer) Snapshot() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Entry, r.count)
	for i := 0; i < r.count; i++ {
		out[i] = r.items[(r.head+i)%r.cap]
	}
	return out
}

// Len returns the current number of entries stored.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Cap returns the maximum capacity of the buffer.
func (r *RingBuffer) Cap() int {
	return r.cap
}
