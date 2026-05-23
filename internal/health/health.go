package health

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status of a source.
type Status struct {
	Source    string    `json:"source"`
	Healthy   bool      `json:"healthy"`
	LastSeen  time.Time `json:"last_seen"`
	ErrorMsg  string    `json:"error,omitempty"`
}

// Checker tracks per-source health and exposes an HTTP handler.
type Checker struct {
	mu       sync.RWMutex
	sources  map[string]*Status
	started  time.Time
}

// New creates a new Checker.
func New() *Checker {
	return &Checker{
		sources: make(map[string]*Status),
		started: time.Now(),
	}
}

// SetHealthy marks a source as healthy with the current timestamp.
func (c *Checker) SetHealthy(source string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sources[source] = &Status{
		Source:   source,
		Healthy:  true,
		LastSeen: time.Now(),
	}
}

// SetUnhealthy marks a source as unhealthy with an error message.
func (c *Checker) SetUnhealthy(source, errMsg string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sources[source] = &Status{
		Source:   source,
		Healthy:  false,
		LastSeen: time.Now(),
		ErrorMsg: errMsg,
	}
}

// Snapshot returns a copy of all source statuses.
func (c *Checker) Snapshot() []Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Status, 0, len(c.sources))
	for _, s := range c.sources {
		out = append(out, *s)
	}
	return out
}

// HTTPHandler returns an http.HandlerFunc that serves health as JSON.
func (c *Checker) HTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := c.Snapshot()
		allHealthy := true
		for _, s := range snap {
			if !s.Healthy {
				allHealthy = false
				break
			}
		}
		payload := map[string]interface{}{
			"ok":      allHealthy,
			"uptime":  time.Since(c.started).String(),
			"sources": snap,
		}
		w.Header().Set("Content-Type", "application/json")
		if !allHealthy {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		_ = json.NewEncoder(w).Encode(payload)
	}
}
