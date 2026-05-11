package processor

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

func makeReader(input string) *reader.Scanner {
	return reader.NewScanner(strings.NewReader(input))
}

func makeWriter() (*output.Writer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	w := output.NewWriter(buf)
	return w, buf
}

func TestPipeline_NoFilters(t *testing.T) {
	sc := makeReader("line1\nline2\nline3\n")
	w, buf := makeWriter()
	p := New(sc, parser.NewParser(), w, nil, nil, nil)
	stats, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Total() != 3 {
		t.Errorf("expected 3 total, got %d", stats.Total())
	}
	if !strings.Contains(buf.String(), "line1") {
		t.Error("expected output to contain 'line1'")
	}
}

func TestPipeline_PatternFilter(t *testing.T) {
	sc := makeReader("error: something\ninfo: ok\nerror: again\n")
	w, buf := makeWriter()
	pf, _ := filter.NewPatternFilter([]string{"error"}, nil)
	p := New(sc, parser.NewParser(), w, pf, nil, nil)
	stats, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Written() != 2 {
		t.Errorf("expected 2 written, got %d", stats.Written())
	}
	if strings.Contains(buf.String(), "info:") {
		t.Error("info line should have been filtered")
	}
}

func TestPipeline_TimeRangeFilter(t *testing.T) {
	start := time.Now().Add(-time.Hour)
	tr := &filter.TimeRange{Start: start}
	sc := makeReader("no-timestamp-line\n")
	w, _ := makeWriter()
	p := New(sc, parser.NewParser(), w, nil, tr, nil)
	_, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPipeline_Deduplicate(t *testing.T) {
	input := "dup-line\ndup-line\nunique\n"
	sc := makeReader(input)
	w, buf := makeWriter()
	dd := NewDeduplicator(100)
	p := New(sc, parser.NewParser(), w, nil, nil, dd)
	stats, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Written() != 2 {
		t.Errorf("expected 2 written after dedup, got %d", stats.Written())
	}
	if strings.Count(buf.String(), "dup-line") != 1 {
		t.Error("expected dup-line to appear exactly once")
	}
}

func TestPipeline_StatsFiltered(t *testing.T) {
	// Verify that stats.Filtered() correctly reflects lines excluded by the pattern filter.
	sc := makeReader("error: one\nwarn: two\nerror: three\ninfo: four\n")
	w, _ := makeWriter()
	pf, _ := filter.NewPatternFilter([]string{"error"}, nil)
	p := New(sc, parser.NewParser(), w, pf, nil, nil)
	stats, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if stats.Total() != 4 {
		t.Errorf("expected 4 total, got %d", stats.Total())
	}
	if stats.Written() != 2 {
		t.Errorf("expected 2 written, got %d", stats.Written())
	}
	if stats.Filtered() != 2 {
		t.Errorf("expected 2 filtered, got %d", stats.Filtered())
	}
}
