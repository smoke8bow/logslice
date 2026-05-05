package filter

import (
	"fmt"
	"time"
)

// TimeRange represents an inclusive start and optional end time window
// used to filter log lines by timestamp.
type TimeRange struct {
	Start time.Time
	End   time.Time
	HasEnd bool
}

// ParseTimeRange parses start and end timestamp strings into a TimeRange.
// Supported layout: "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02".
func ParseTimeRange(start, end string) (*TimeRange, error) {
	if start == "" {
		return nil, fmt.Errorf("start time must not be empty")
	}

	s, err := parseTimestamp(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time %q: %w", start, err)
	}

	tr := &TimeRange{Start: s}

	if end != "" {
		e, err := parseTimestamp(end)
		if err != nil {
			return nil, fmt.Errorf("invalid end time %q: %w", end, err)
		}
		if e.Before(s) {
			return nil, fmt.Errorf("end time %q must not be before start time %q", end, start)
		}
		tr.End = e
		tr.HasEnd = true
	}

	return tr, nil
}

// Contains reports whether t falls within the time range.
func (tr *TimeRange) Contains(t time.Time) bool {
	if t.Before(tr.Start) {
		return false
	}
	if tr.HasEnd && t.After(tr.End) {
		return false
	}
	return true
}

var timestampLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parseTimestamp(s string) (time.Time, error) {
	for _, layout := range timestampLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp format")
}
