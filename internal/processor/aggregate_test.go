package processor

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeAggLine(fields map[string]string) parser.LogLine {
	return parser.LogLine{Fields: fields, Raw: ""}
}

func TestNewAggregator_EmptyGroupField(t *testing.T) {
	_, err := NewAggregator("", "", AggCount)
	if err == nil {
		t.Fatal("expected error for empty groupField")
	}
}

func TestNewAggregator_UnknownMode(t *testing.T) {
	_, err := NewAggregator("level", "", "median")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewAggregator_MissingValueFieldForSum(t *testing.T) {
	_, err := NewAggregator("level", "", AggSum)
	if err == nil {
		t.Fatal("expected error when valueField empty for sum mode")
	}
}

func TestAggregator_CountMode(t *testing.T) {
	a, err := NewAggregator("level", "", AggCount)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.Ingest(makeAggLine(map[string]string{"level": "info"}))
	a.Ingest(makeAggLine(map[string]string{"level": "info"}))
	a.Ingest(makeAggLine(map[string]string{"level": "error"}))
	res := a.Results()
	if len(res) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(res))
	}
	for _, r := range res {
		switch r.Key {
		case "info":
			if r.Value != 2 {
				t.Errorf("info count: want 2, got %v", r.Value)
			}
		case "error":
			if r.Value != 1 {
				t.Errorf("error count: want 1, got %v", r.Value)
			}
		}
	}
}

func TestAggregator_SumMode(t *testing.T) {
	a, _ := NewAggregator("service", "latency", AggSum)
	a.Ingest(makeAggLine(map[string]string{"service": "api", "latency": "10.5"}))
	a.Ingest(makeAggLine(map[string]string{"service": "api", "latency": "4.5"}))
	a.Ingest(makeAggLine(map[string]string{"service": "db", "latency": "20"}))
	res := a.Results()
	for _, r := range res {
		if r.Key == "api" && r.Value != 15.0 {
			t.Errorf("api sum: want 15, got %v", r.Value)
		}
		if r.Key == "db" && r.Value != 20.0 {
			t.Errorf("db sum: want 20, got %v", r.Value)
		}
	}
}

func TestAggregator_MinMaxMode(t *testing.T) {
	amin, _ := NewAggregator("host", "cpu", AggMin)
	amax, _ := NewAggregator("host", "cpu", AggMax)
	lines := []map[string]string{
		{"host": "h1", "cpu": "30"},
		{"host": "h1", "cpu": "80"},
		{"host": "h1", "cpu": "50"},
	}
	for _, f := range lines {
		amin.Ingest(makeAggLine(f))
		amax.Ingest(makeAggLine(f))
	}
	for _, r := range amin.Results() {
		if r.Key == "h1" && r.Value != 30 {
			t.Errorf("min: want 30, got %v", r.Value)
		}
	}
	for _, r := range amax.Results() {
		if r.Key == "h1" && r.Value != 80 {
			t.Errorf("max: want 80, got %v", r.Value)
		}
	}
}

func TestAggregator_MissingGroupField_UsesNone(t *testing.T) {
	a, _ := NewAggregator("level", "", AggCount)
	a.Ingest(makeAggLine(map[string]string{"msg": "hello"}))
	res := a.Results()
	if len(res) != 1 || res[0].Key != "(none)" {
		t.Errorf("expected (none) bucket, got %+v", res)
	}
}

func TestAggregator_Reset(t *testing.T) {
	a, _ := NewAggregator("level", "", AggCount)
	a.Ingest(makeAggLine(map[string]string{"level": "info"}))
	a.Reset()
	if len(a.Results()) != 0 {
		t.Error("expected empty results after reset")
	}
}
