package sampling

import (
	"testing"

	"github.com/your/logpipe/internal/normalize"
)

func entry(source string) normalize.Entry {
	return normalize.Entry{Source: source, Message: "test"}
}

func TestKeep_RateOne_AlwaysKeeps(t *testing.T) {
	s := New(Config{Strategy: StrategyRandom, Rate: 1.0})
	for i := 0; i < 100; i++ {
		if !s.Keep(entry("svc")) {
			t.Fatal("expected entry to be kept with rate=1.0")
		}
	}
}

func TestKeep_RateZero_AlwaysDrops(t *testing.T) {
	s := New(Config{Strategy: StrategyRandom, Rate: 0.0})
	for i := 0; i < 100; i++ {
		if s.Keep(entry("svc")) {
			t.Fatal("expected entry to be dropped with rate=0.0")
		}
	}
}

func TestKeep_RateClamped(t *testing.T) {
	s := New(Config{Rate: 1.5})
	if s.cfg.Rate != 1.0 {
		t.Fatalf("expected rate clamped to 1.0, got %f", s.cfg.Rate)
	}
	s2 := New(Config{Rate: -0.5})
	if s2.cfg.Rate != 0.0 {
		t.Fatalf("expected rate clamped to 0.0, got %f", s2.cfg.Rate)
	}
}

func TestKeep_PerSourceOverride(t *testing.T) {
	s := New(Config{
		Rate: 1.0,
		PerSource: map[string]float64{"noisy": 0.0},
	})
	if s.Keep(entry("noisy")) {
		t.Fatal("expected noisy source to be dropped")
	}
	if !s.Keep(entry("other")) {
		t.Fatal("expected other source to be kept")
	}
}

func TestKeep_StatisticalRate(t *testing.T) {
	s := New(Config{Rate: 0.5})
	kept := 0
	const n = 10000
	for i := 0; i < n; i++ {
		if s.Keep(entry("svc")) {
			kept++
		}
	}
	ratio := float64(kept) / n
	if ratio < 0.40 || ratio > 0.60 {
		t.Fatalf("expected ~50%% kept, got %.2f%%", ratio*100)
	}
}
