package tail

import (
	"context"
	"fmt"
	"sync"
)

// Source describes a log source to tail.
type Source struct {
	Name string
	Path string
}

// Manager supervises multiple Tailers and fans their output into a single channel.
type Manager struct {
	sources []Source
	out     chan Line
}

// NewManager creates a Manager for the given sources.
func NewManager(sources []Source) *Manager {
	return &Manager{
		sources: sources,
		out:     make(chan Line, 256),
	}
}

// Lines returns the read-only channel of aggregated log lines.
func (m *Manager) Lines() <-chan Line {
	return m.out
}

// Run starts a Tailer for each source and blocks until ctx is cancelled.
// Errors from individual tailers are printed but do not stop the manager.
func (m *Manager) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, src := range m.sources {
		src := src
		wg.Add(1)
		go func() {
			defer wg.Done()
			t := New(src.Path, src.Name, m.out)
			if err := t.Run(ctx); err != nil && ctx.Err() == nil {
				fmt.Printf("tailer error [%s]: %v\n", src.Name, err)
			}
		}()
	}

	wg.Wait()
	close(m.out)
	return nil
}
