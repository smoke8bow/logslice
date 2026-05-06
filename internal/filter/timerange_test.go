package filter_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
)

func TestParseTimeRange_ValidStartOnly(t *testing.T) {
	tr, err := filter.ParseTimeRange("2024-01-01T10:00:00Z", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.End != (time.Time{}) {
		t.Errorf("expected zero End, got %v", tr.End)
	}
}

func TestParseTimeRange_ValidStartAndEnd(t *testing.T) {
	tr, err := filter.ParseTimeRange("2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.End.IsZero() {
		t.Error("expected non-zero End")
	}
}

func TestParseTimeRange_EndBeforeStart(t *testing.T) {
	_, err := filter.ParseTimeRange("2024-01-01T12:00:00Z", "2024-01-01T10:00:00Z")
	if err == nil {
		t.Error("expected error for end before start")
	}
}

func TestParseTimeRange_EmptyStart(t *testing.T) {
	_, err := filter.ParseTimeRange("", "2024-01-01T12:00:00Z")
	if err == nil {
		t.Error("expected error for empty start")
	}
}

func TestParseTimeRange_InvalidFormat(t *testing.T) {
	_, err := filter.ParseTimeRange("not-a-date", "")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestTimeRange_Contains(t *testing.T) {
	tr, _ := filter.ParseTimeRange("2024-01-01T10:00:00Z", "2024-01-01T12:00:00Z")
	cases := []struct {
		ts   string
		want bool
	}{
		{"2024-01-01T09:59:59Z", false},
		{"2024-01-01T10:00:00Z", true},
		{"2024-01-01T11:00:00Z", true},
		{"2024-01-01T12:00:00Z", false},
		{"2024-01-01T13:00:00Z", false},
	}
	for _, c := range cases {
		t0, _ := time.Parse(time.RFC3339, c.ts)
		if got := tr.Contains(t0); got != c.want {
			t.Errorf("Contains(%s) = %v, want %v", c.ts, got, c.want)
		}
	}
}
