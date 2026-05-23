package sampling

import (
	"math/rand"
	"sync"
	"time"

	"github.com/your/logpipe/internal/normalize"
)

// Strategy defines how sampling is applied.
type Strategy string

const (
	StrategyRandom     Strategy = "random"
	StrategyRateBased  Strategy = "rate_based"
)

// Config holds sampler configuration.
type Config struct {
	Strategy   Strategy
	Rate       float64 // 0.0 to 1.0, fraction of entries to keep
	PerSource  map[string]float64 // optional per-source override
}

// Sampler decides whether a log entry should be kept.
type Sampler struct {
	cfg  Config
	rng  *rand.Rand
	mu   sync.Mutex
}

// New creates a Sampler from cfg. Rate is clamped to [0.0, 1.0].
func New(cfg Config) *Sampler {
	if cfg.Rate < 0 {
		cfg.Rate = 0
	}
	if cfg.Rate > 1 {
		cfg.Rate = 1
	}
	return &Sampler{
		cfg: cfg,
		//nolint:gosec
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Keep returns true if the entry should be kept according to the sampling policy.
func (s *Sampler) Keep(entry normalize.Entry) bool {
	rate := s.cfg.Rate
	if override, ok := s.cfg.PerSource[entry.Source]; ok {
		rate = override
	}
	if rate >= 1.0 {
		return true
	}
	if rate <= 0.0 {
		return false
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < rate
}
