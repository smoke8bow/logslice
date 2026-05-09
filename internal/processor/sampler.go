package processor

import (
	"math/rand"
	"sync/atomic"
)

// Sampler provides line-level sampling for log processing pipelines.
// It supports both deterministic (every Nth line) and random sampling modes.
type Sampler struct {
	mode     string
	rate     float64
	nth      uint64
	counter  atomic.Uint64
	rng      *rand.Rand
}

// NewSampler creates a Sampler. mode must be "rate" or "nth".
// For "rate" mode, rate must be in (0.0, 1.0].
// For "nth" mode, nth must be >= 1 (keep every nth line).
func NewSampler(mode string, rate float64, nth uint64) (*Sampler, error) {
	switch mode {
	case "rate":
		if rate <= 0.0 || rate > 1.0 {
			return nil, fmt.Errorf("sampler: rate must be in (0.0, 1.0], got %f", rate)
		}
	case "nth":
		if nth < 1 {
			return nil, fmt.Errorf("sampler: nth must be >= 1, got %d", nth)
		}
	default:
		return nil, fmt.Errorf("sampler: unknown mode %q, must be \"rate\" or \"nth\"", mode)
	}
	return &Sampler{
		mode: mode,
		rate: rate,
		nth:  nth,
		rng:  rand.New(rand.NewSource(rand.Int63())),
	}, nil
}

// Keep returns true if the current line should be kept based on the sampling strategy.
func (s *Sampler) Keep() bool {
	count := s.counter.Add(1)
	switch s.mode {
	case "rate":
		return s.rng.Float64() < s.rate
	case "nth":
		return count%s.nth == 0
	}
	return true
}

// Reset resets the internal line counter (useful between files or test runs).
func (s *Sampler) Reset() {
	s.counter.Store(0)
}
