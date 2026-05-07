package processor

import (
	"fmt"
	"io"
	"time"
)

// Stats holds counters collected during a pipeline or parallel run.
type Stats struct {
	LinesRead    int64
	LinesWritten int64
	LinesDropped int64
	Duration     time.Duration
}

// PassRate returns the percentage of lines that passed all filters.
func (s *Stats) PassRate() float64 {
	if s.LinesRead == 0 {
		return 0
	}
	return float64(s.LinesWritten) / float64(s.LinesRead) * 100
}

// Print writes a human-readable summary of the stats to w.
func (s *Stats) Print(w io.Writer) {
	fmt.Fprintf(w, "Lines read:    %d\n", s.LinesRead)
	fmt.Fprintf(w, "Lines written: %d\n", s.LinesWritten)
	fmt.Fprintf(w, "Lines dropped: %d\n", s.LinesDropped)
	fmt.Fprintf(w, "Pass rate:     %.1f%%\n", s.PassRate())
	fmt.Fprintf(w, "Duration:      %s\n", s.Duration.Round(time.Millisecond))
}

// Merge combines another Stats into this one (for aggregating worker results).
func (s *Stats) Merge(other Stats) {
	s.LinesRead += other.LinesRead
	s.LinesWritten += other.LinesWritten
	s.LinesDropped += other.LinesDropped
	if other.Duration > s.Duration {
		s.Duration = other.Duration
	}
}
