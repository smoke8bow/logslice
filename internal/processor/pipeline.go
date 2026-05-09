package processor

import (
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

// Pipeline processes log lines from a reader through optional filters and writes results.
type Pipeline struct {
	scanner   *reader.Scanner
	parser    *parser.Parser
	writer    *output.Writer
	pattern   *filter.PatternFilter
	timeRange *filter.TimeRange
	dedupe    *Deduplicator
	stats     *Stats
}

// New constructs a Pipeline with the given components.
func New(
	scanner *reader.Scanner,
	prs *parser.Parser,
	w *output.Writer,
	pf *filter.PatternFilter,
	tr *filter.TimeRange,
	dd *Deduplicator,
) *Pipeline {
	return &Pipeline{
		scanner:   scanner,
		parser:    prs,
		writer:    w,
		pattern:   pf,
		timeRange: tr,
		dedupe:    dd,
		stats:     NewStats(),
	}
}

// Run executes the pipeline and returns the final stats.
func (p *Pipeline) Run() (*Stats, error) {
	for p.scanner.Scan() {
		line := p.scanner.Text()
		p.stats.Inc()

		if p.dedupe != nil && p.dedupe.IsDuplicate(line) {
			p.stats.Skip()
			continue
		}

		parsed := p.parser.Parse(line)

		if p.timeRange != nil && !parsed.Time.IsZero() {
			if !p.timeRange.Contains(parsed.Time) {
				p.stats.Skip()
				continue
			}
		}

		if p.pattern != nil && !p.pattern.Match(line) {
			p.stats.Skip()
			continue
		}

		if err := p.writer.WriteLine(line); err != nil {
			return nil, err
		}
		p.stats.Write()
	}

	if err := p.scanner.Err(); err != nil {
		return nil, err
	}

	p.stats.Finish()
	return p.stats, nil
}
