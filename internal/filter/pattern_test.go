package filter

import (
	"testing"
)

func TestNewPatternFilter_InvalidInclude(t *testing.T) {
	_, err := NewPatternFilter([]string{"[invalid"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid include pattern, got nil")
	}
}

func TestNewPatternFilter_InvalidExclude(t *testing.T) {
	_, err := NewPatternFilter(nil, []string{"(unclosed"})
	if err == nil {
		t.Fatal("expected error for invalid exclude pattern, got nil")
	}
}

func TestPatternFilter_NoFilters(t *testing.T) {
	pf, err := NewPatternFilter(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pf.Match("any line at all") {
		t.Error("expected match with no filters set")
	}
	if pf.HasFilters() {
		t.Error("expected HasFilters to return false")
	}
}

func TestPatternFilter_IncludeOnly(t *testing.T) {
	pf, err := NewPatternFilter([]string{"ERROR", "WARN"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pf.Match("2024-01-01 ERROR something went wrong") {
		t.Error("expected match for ERROR line")
	}
	if !pf.Match("2024-01-01 WARN low disk space") {
		t.Error("expected match for WARN line")
	}
	if pf.Match("2024-01-01 INFO all good") {
		t.Error("expected no match for INFO line")
	}
}

func TestPatternFilter_ExcludeOnly(t *testing.T) {
	pf, err := NewPatternFilter(nil, []string{"DEBUG"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pf.Match("2024-01-01 DEBUG verbose output") {
		t.Error("expected no match for excluded DEBUG line")
	}
	if !pf.Match("2024-01-01 INFO startup complete") {
		t.Error("expected match for non-excluded INFO line")
	}
}

func TestPatternFilter_IncludeAndExclude(t *testing.T) {
	pf, err := NewPatternFilter([]string{"ERROR"}, []string{"timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !pf.Match("ERROR: disk full") {
		t.Error("expected match for ERROR without excluded term")
	}
	if pf.Match("ERROR: connection timeout") {
		t.Error("expected no match for ERROR with excluded term 'timeout'")
	}
	if pf.Match("INFO: all systems nominal") {
		t.Error("expected no match for INFO line not in includes")
	}
}

func TestPatternFilter_EmptyStringPatterns(t *testing.T) {
	pf, err := NewPatternFilter([]string{"  ", ""}, []string{""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pf.HasFilters() {
		t.Error("expected HasFilters false when only blank patterns provided")
	}
}
