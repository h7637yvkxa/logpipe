package dedupe

import "time"

// Pipeline wraps a Deduplicator and filters a channel of log entries,
// forwarding only non-duplicate entries to the output channel.
type Pipeline struct {
	deduper *Deduplicator
}

// NewPipeline creates a Pipeline using the given deduplication window.
func NewPipeline(window time.Duration) *Pipeline {
	return &Pipeline{deduper: New(window)}
}

// Run reads entries from in, suppresses duplicates, and writes unique entries
// to the returned channel. The output channel is closed when in is closed.
func (p *Pipeline) Run(in <-chan Entry) <-chan Entry {
	out := make(chan Entry, cap(in))
	go func() {
		defer close(out)
		for entry := range in {
			if !p.deduper.IsDuplicate(entry) {
				out <- entry
			}
		}
	}()
	return out
}
