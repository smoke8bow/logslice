package processor

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Stats tracks processing metrics for a pipeline run.
type Stats struct {
	LinesRead     atomic.Int64
	LinesMatched  atomic.Int64
	LinesDropped  atomic.Int64
	BytesRead     atomic.Int64
	StartTime     time.Time
	EndTime       time.Time
}

// NewStats creates a new Stats instance with the start time set.
func NewStats() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// Finish marks the end time of processing.
func (s *Stats) Finish() {
	s.EndTime = time.Now()
}

// Duration returns the elapsed processing time.
func (s *Stats) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// Summary returns a human-readable summary string.
func (s *Stats) Summary() string {
	return fmt.Sprintf(
		"lines_read=%d lines_matched=%d lines_dropped=%d bytes_read=%d duration=%s",
		s.LinesRead.Load(),
		s.LinesMatched.Load(),
		s.LinesDropped.Load(),
		s.BytesRead.Load(),
		s.Duration().Round(time.Millisecond),
	)
}
