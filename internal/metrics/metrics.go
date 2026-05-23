package metrics

import (
	"sync"
	"sync/atomic"
)

// Counters holds runtime statistics for logpipe.
type Counters struct {
	LinesRead      atomic.Int64
	LinesFiltered  atomic.Int64
	LinesWritten   atomic.Int64
	ParseErrors    atomic.Int64
	mu             sync.RWMutex
	sourceCounters map[string]*atomic.Int64
}

// New returns an initialised Counters instance.
func New() *Counters {
	return &Counters{
		sourceCounters: make(map[string]*atomic.Int64),
	}
}

// IncRead increments the global lines-read counter and the per-source counter.
func (c *Counters) IncRead(source string) {
	c.LinesRead.Add(1)
	c.sourceCounter(source).Add(1)
}

// IncFiltered increments the lines-filtered counter.
func (c *Counters) IncFiltered() {
	c.LinesFiltered.Add(1)
}

// IncWritten increments the lines-written counter.
func (c *Counters) IncWritten() {
	c.LinesWritten.Add(1)
}

// IncParseError increments the parse-error counter.
func (c *Counters) IncParseError() {
	c.ParseErrors.Add(1)
}

// Snapshot returns a point-in-time copy of all counters.
func (c *Counters) Snapshot() map[string]int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snap := map[string]int64{
		"lines_read":     c.LinesRead.Load(),
		"lines_filtered": c.LinesFiltered.Load(),
		"lines_written":  c.LinesWritten.Load(),
		"parse_errors":   c.ParseErrors.Load(),
	}
	for src, ctr := range c.sourceCounters {
		snap["source."+src] = ctr.Load()
	}
	return snap
}

// sourceCounter returns (creating if necessary) the per-source atomic counter.
func (c *Counters) sourceCounter(source string) *atomic.Int64 {
	c.mu.RLock()
	if ctr, ok := c.sourceCounters[source]; ok {
		c.mu.RUnlock()
		return ctr
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if ctr, ok := c.sourceCounters[source]; ok {
		return ctr
	}
	ctr := &atomic.Int64{}
	c.sourceCounters[source] = ctr
	return ctr
}
