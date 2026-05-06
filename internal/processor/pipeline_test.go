package processor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/processor"
	"github.com/logslice/logslice/internal/reader"
)

func makeReader(t *testing.T, content string) *reader.LineReader {
	t.Helper()
	r, err := reader.NewLineReader(strings.NewReader(content))
	if err != nil {
		t.Fatalf("NewLineReader: %v", err)
	}
	return r
}

func makeWriter(buf *bytes.Buffer) *output.Writer {
	return output.NewWriter(buf)
}

func TestPipeline_NoFilters(t *testing.T) {
	buf := &bytes.Buffer{}
	p := processor.New(makeReader(t, "line1\nline2\nline3\n"), processor.Config{
		Writer: makeWriter(buf),
		Format: output.FormatRaw,
	})
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3 lines written, got %d", n)
	}
}

func TestPipeline_PatternFilter(t *testing.T) {
	buf := &bytes.Buffer{}
	pf, _ := filter.NewPatternFilter([]string{"error"}, nil)
	p := processor.New(makeReader(t, "info: ok\nerror: bad\nwarn: meh\nerror: also bad\n"), processor.Config{
		Pattern: pf,
		Writer:  makeWriter(buf),
		Format:  output.FormatRaw,
	})
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 lines, got %d", n)
	}
}

func TestPipeline_TimeRangeFilter(t *testing.T) {
	buf := &bytes.Buffer{}
	start := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := &filter.TimeRange{Start: start, End: end}

	lines := "before\ninside\nafter\n"
	times := map[string]time.Time{
		"before": time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		"inside": time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		"after":  time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
	}
	tsFn := func(line string) (time.Time, bool) {
		if t, ok := times[line]; ok {
			return t, true
		}
		return time.Time{}, false
	}

	p := processor.New(makeReader(t, lines), processor.Config{
		TimeRange:   tr,
		Writer:      makeWriter(buf),
		Format:      output.FormatRaw,
		TimestampFn: tsFn,
	})
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 line, got %d", n)
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	buf := &bytes.Buffer{}
	p := processor.New(makeReader(t, ""), processor.Config{
		Writer: makeWriter(buf),
		Format: output.FormatRaw,
	})
	n, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 lines, got %d", n)
	}
}
