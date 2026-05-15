package processor

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/reader"
)

// AggregateConfig holds configuration for an aggregation run.
type AggregateConfig struct {
	InputFile   string
	GroupField  string
	ValueField  string
	Mode        AggregateMode
	OutputFmt   output.Format
	Writer      io.Writer
}

// RunAggregate reads all lines from InputFile, feeds them into an Aggregator,
// and writes the results to the configured writer.
func RunAggregate(cfg AggregateConfig) error {
	agg, err := NewAggregator(cfg.GroupField, cfg.ValueField, cfg.Mode)
	if err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}

	lr, err := reader.NewLineReader(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("aggregate: open input: %w", err)
	}
	defer lr.Close()

	p := parser.NewParser()
	for {
		line, readErr := lr.ReadLine()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return fmt.Errorf("aggregate: read: %w", readErr)
		}
		parsed := p.Parse(line)
		agg.Ingest(parsed)
	}

	w, err := output.NewWriter(cfg.Writer, cfg.OutputFmt)
	if err != nil {
		return fmt.Errorf("aggregate: writer: %w", err)
	}

	results := agg.Results()
	for _, r := range results {
		var rendered string
		switch cfg.OutputFmt {
		case output.FormatJSON:
			if cfg.Mode == AggCount {
				rendered = fmt.Sprintf(`{"group":%q,"count":%d}`, r.Key, r.Count)
			} else {
				rendered = fmt.Sprintf(`{"group":%q,"%s":%g,"count":%d}`,
					r.Key, strings.ToLower(string(cfg.Mode)), r.Value, r.Count)
			}
		default:
			if cfg.Mode == AggCount {
				rendered = fmt.Sprintf("%s\t%d", r.Key, r.Count)
			} else {
				rendered = fmt.Sprintf("%s\t%g\t(n=%d)", r.Key, r.Value, r.Count)
			}
		}
		if wErr := w.WriteLine(rendered); wErr != nil {
			return fmt.Errorf("aggregate: write: %w", wErr)
		}
	}
	return w.Flush()
}
