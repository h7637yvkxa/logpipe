package transform

import "github.com/user/logpipe/internal/normalize"

// Pipeline reads entries from in, applies the Transformer, and sends
// results to the returned channel. The output channel is closed when
// in is closed.
func Pipeline(in <-chan normalize.Entry, t *Transformer) <-chan normalize.Entry {
	out := make(chan normalize.Entry)
	go func() {
		defer close(out)
		for entry := range in {
			out <- t.Apply(entry)
		}
	}()
	return out
}
