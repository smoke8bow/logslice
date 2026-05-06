package filter

import (
	"errors"
	"fmt"
	"time"
)

// TimeRange represents an inclusive start / exclusive-end time window.
// End is optional; a zero End means "no upper bound".
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Contains reports whether t falls within the time range.
func (tr *TimeRange) Contains(t time.Time) bool {
	if t.Before(tr.Start) {
		return false
	}
	if !tr.End.IsZero() && !t.Before(tr.End) {
		return false
	}
	return true
}

// supportedLayouts lists the timestamp formats tried in order.
var supportedLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// parseTimestamp attempts to parse s using each supported layout.
func parseTimestamp(s string) (time.Time, error) {
	for _, layout := range supportedLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp format: %q", s)
}

// ParseTimeRange builds a TimeRange from string representations of start and
// (optionally) end timestamps. start must not be empty.
func ParseTimeRange(start, end string) (*TimeRange, error) {
	if start == "" {
		return nil, errors.New("start timestamp is required")
	}
	s, err := parseTimestamp(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start: %w", err)
	}
	tr := &TimeRange{Start: s}
	if end != "" {
		e, err := parseTimestamp(end)
		if err != nil {
			return nil, fmt.Errorf("invalid end: %w", err)
		}
		if !e.After(s) {
			return nil, errors.New("end timestamp must be after start timestamp")
		}
		tr.End = e
	}
	return tr, nil
}
