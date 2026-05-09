package processor

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Stats tracks pipeline processing counters and timing.
type Stats struct {
	total   int64
	skipped int64
	written int64
	start   time.Time
	end     time.Time
}

// NewStats initialises a Stats instance with the current time.
func NewStats() *Stats {
	return &Stats{start: time.Now()}
}

// Inc increments the total lines seen counter.
func (s *Stats) Inc() {
	atomic.AddInt64(&s.total, 1)
}

// Skip increments the skipped lines counter.
func (s *Stats) Skip() {
	atomic.AddInt64(&s.skipped, 1)
}

// Write increments the written lines counter.
func (s *Stats) Write() {
	atomic.AddInt64(&s.written, 1)
}

// Finish records the end time.
func (s *Stats) Finish() {
	s.end = time.Now()
}

// Total returns total lines processed.
func (s *Stats) Total() int64 { return atomic.LoadInt64(&s.total) }

// Skipped returns total lines skipped.
func (s *Stats) Skipped() int64 { return atomic.LoadInt64(&s.skipped) }

// Written returns total lines written.
func (s *Stats) Written() int64 { return atomic.LoadInt64(&s.written) }

// Duration returns elapsed time. If Finish has not been called, returns time since start.
func (s *Stats) Duration() time.Duration {
	if s.end.IsZero() {
		return time.Since(s.start)
	}
	return s.end.Sub(s.start)
}

// Summary returns a human-readable summary string.
func (s *Stats) Summary() string {
	return fmt.Sprintf(
		"total=%d written=%d skipped=%d duration=%s",
		s.Total(), s.Written(), s.Skipped(), s.Duration().Round(time.Millisecond),
	)
}
