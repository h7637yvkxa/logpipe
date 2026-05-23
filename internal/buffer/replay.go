package buffer

import (
	"encoding/json"
	"io"
)

// ReplayOptions controls how buffered entries are replayed.
type ReplayOptions struct {
	// Source filters replay to a specific source label; empty means all.
	Source string
	// Limit caps the number of entries replayed; 0 means all.
	Limit int
}

// Replay writes buffered entries from rb to w as newline-delimited JSON.
// Each written object contains "source" and "line" fields.
// Entries are filtered and limited according to opts.
func Replay(rb *RingBuffer, w io.Writer, opts ReplayOptions) error {
	snap := rb.Snapshot()

	var filtered []Entry
	for _, e := range snap {
		if opts.Source != "" && e.Source != opts.Source {
			continue
		}
		filtered = append(filtered, e)
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[len(filtered)-opts.Limit:]
	}

	enc := json.NewEncoder(w)
	for _, e := range filtered {
		if err := enc.Encode(map[string]string{
			"source": e.Source,
			"line":   e.Line,
		}); err != nil {
			return err
		}
	}
	return nil
}
