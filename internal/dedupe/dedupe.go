package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

// Entry represents a normalized log entry for deduplication purposes.
type Entry map[string]interface{}

// Deduplicator tracks recently seen log entries and suppresses duplicates
// within a configurable time window.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// New creates a Deduplicator with the given deduplication window.
// Entries with identical content seen within the window are considered duplicates.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// IsDuplicate returns true if an identical entry was seen within the window.
// It also records the entry if it is not a duplicate.
func (d *Deduplicator) IsDuplicate(entry Entry) bool {
	h := hash(entry)
	if h == "" {
		return false
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	d.evict(now)

	if _, seen := d.seen[h]; seen {
		return true
	}

	d.seen[h] = now
	return false
}

// evict removes entries outside the deduplication window. Must be called with mu held.
func (d *Deduplicator) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) > d.window {
			delete(d.seen, k)
		}
	}
}

// hash returns a stable SHA-256 hex digest of the entry's JSON representation.
func hash(entry Entry) string {
	b, err := json.Marshal(entry)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
