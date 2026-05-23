package sampling

import "github.com/your/logpipe/internal/normalize"

// Pipeline reads from in, drops entries that do not pass the sampler,
// and forwards the rest to the returned channel.
func Pipeline(in <-chan normalize.Entry, s *Sampler) <-chan normalize.Entry {
	out := make(chan normalize.Entry, cap(in))
	go func() {
		defer close(out)
		for entry := range in {
			if s.Keep(entry) {
				out <- entry
			}
		}
	}()
	return out
}
