package metrics

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"
)

// Reporter periodically writes a JSON metrics snapshot to a writer.
type Reporter struct {
	counters  *Counters
	interval  time.Duration
	out       io.Writer
}

// NewReporter creates a Reporter that writes to out every interval.
// If out is nil, os.Stderr is used.
func NewReporter(c *Counters, interval time.Duration, out io.Writer) *Reporter {
	if out == nil {
		out = os.Stderr
	}
	return &Reporter{
		counters: c,
		interval: interval,
		out:      out,
	}
}

// Run starts the periodic reporting loop. It blocks until ctx is cancelled.
func (r *Reporter) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.write() // final flush
			return
		case <-ticker.C:
			r.write()
		}
	}
}

func (r *Reporter) write() {
	snap := r.counters.Snapshot()
	snap["_ts"] = time.Now().Unix()
	data, err := json.Marshal(snap)
	if err != nil {
		return
	}
	data = append(data, '\n')
	_, _ = r.out.Write(data)
}
