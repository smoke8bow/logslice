package processor

import (
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/reader"
)

// Pipeline ties together a line reader, time range filter, pattern filter,
// and output writer into a single processing pass.
type Pipeline struct {
	reader      *reader.LineReader
	timeRange   *filter.TimeRange
	pattern     *filter.PatternFilter
	writer      *output.Writer
	format      output.Format
	timestampFn func(string) (time.Time, bool)
}

// Config holds all options needed to construct a Pipeline.
type Config struct {
	TimeRange   *filter.TimeRange
	Pattern     *filter.PatternFilter
	Writer      *output.Writer
	Format      output.Format
	TimestampFn func(string) (time.Time, bool)
}

// New creates a Pipeline from a LineReader and a Config.
func New(r *reader.LineReader, cfg Config) *Pipeline {
	ts := cfg.TimestampFn
	if ts == nil {
		ts = func(string) (time.Time, bool) { return time.Time{}, false }
	}
	return &Pipeline{
		reader:      r,
		timeRange:   cfg.TimeRange,
		pattern:     cfg.Pattern,
		writer:      cfg.Writer,
		format:      cfg.Format,
		timestampFn: ts,
	}
}

// Run iterates over every line and applies filters, writing matching lines.
// Returns the number of lines written and any write error.
func (p *Pipeline) Run() (int, error) {
	written := 0
	for {
		line, err := p.reader.ReadLine()
		if err != nil {
			break
		}
		if !p.accept(line) {
			continue
		}
		formatted := output.FormatLine(p.format, line)
		if werr := p.writer.WriteLine(formatted); werr != nil {
			return written, werr
		}
		written++
	}
	return written, nil
}

// accept returns true when the line passes all active filters.
func (p *Pipeline) accept(line string) bool {
	if p.timeRange != nil {
		if t, ok := p.timestampFn(line); ok {
			if !p.timeRange.Contains(t) {
				return false
			}
		}
	}
	if p.pattern != nil {
		if !p.pattern.Match(line) {
			return false
		}
	}
	return true
}
