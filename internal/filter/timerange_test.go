package filter

import (
	"testing"
	"time"
)

func TestParseTimeRange_ValidStartOnly(t *testing.T) {
	tr, err := ParseTimeRange("2024-03-01T10:00:00", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.HasEnd {
		t.Error("expected HasEnd to be false")
	}
	expected := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	if !tr.Start.Equal(expected) {
		t.Errorf("Start = %v, want %v", tr.Start, expected)
	}
}

func TestParseTimeRange_ValidStartAndEnd(t *testing.T) {
	tr, err := ParseTimeRange("2024-03-01", "2024-03-02")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.HasEnd {
		t.Error("expected HasEnd to be true")
	}
}

func TestParseTimeRange_EndBeforeStart(t *testing.T) {
	_, err := ParseTimeRange("2024-03-05T00:00:00", "2024-03-01T00:00:00")
	if err == nil {
		t.Fatal("expected error when end is before start")
	}
}

func TestParseTimeRange_EmptyStart(t *testing.T) {
	_, err := ParseTimeRange("", "2024-03-01")
	if err == nil {
		t.Fatal("expected error for empty start")
	}
}

func TestParseTimeRange_InvalidFormat(t *testing.T) {
	_, err := ParseTimeRange("not-a-date", "")
	if err == nil {
		t.Fatal("expected error for invalid timestamp format")
	}
}

func TestTimeRange_Contains(t *testing.T) {
	tr, _ := ParseTimeRange("2024-03-01T08:00:00", "2024-03-01T18:00:00")

	cases := []struct {
		ts   string
		want bool
	}{
		{"2024-03-01T07:59:59", false},
		{"2024-03-01T08:00:00", true},
		{"2024-03-01T12:00:00", true},
		{"2024-03-01T18:00:00", true},
		{"2024-03-01T18:00:01", false},
	}

	for _, c := range cases {
		t, err := time.Parse("2006-01-02T15:04:05", c.ts)
		if err != nil {
			test.Fatalf("bad test timestamp %q: %v", c.ts, err)
		}
		got := tr.Contains(t)
		if got != c.want {
			test.Errorf("Contains(%q) = %v, want %v", c.ts, got, c.want)
		}
	}
}
