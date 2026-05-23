// Package sampling provides probabilistic log entry sampling for logpipe.
//
// A Sampler is configured with a global keep-rate (0.0–1.0) and optional
// per-source overrides. The Pipeline helper wraps a Sampler around an entry
// channel so it can be inserted into the processing pipeline between
// normalization and filtering.
//
// Example usage:
//
//	sampler := sampling.New(sampling.Config{
//		Strategy:  sampling.StrategyRandom,
//		Rate:      0.25, // keep 25 % of all entries
//		PerSource: map[string]float64{"debug-svc": 0.05},
//	})
//	sampled := sampling.Pipeline(normalizedCh, sampler)
package sampling
