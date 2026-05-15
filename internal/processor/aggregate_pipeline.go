package processor

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/parser"
)

// AggregateStage is a pipeline stage that accumulates lines into an Aggregator
// and passes them through unchanged. Call Flush to emit results after the
// pipeline has drained.
type AggregateStage struct {
	agg    *Aggregator
	result []AggregateResult
}

// NewAggregateStage creates a pipeline-compatible aggregation stage.
func NewAggregateStage(groupField, valueField string, mode AggregateMode) (*AggregateStage, error) {
	agg, err := NewAggregator(groupField, valueField, mode)
	if err != nil {
		return nil, err
	}
	return &AggregateStage{agg: agg}, nil
}

// Process ingests the line and returns it unchanged so downstream stages
// still receive every line.
func (s *AggregateStage) Process(line parser.LogLine) (parser.LogLine, bool) {
	s.agg.Ingest(line)
	return line, true
}

// Flush finalises the aggregation and caches results. It is idempotent.
func (s *AggregateStage) Flush() {
	if s.result == nil {
		s.result = s.agg.Results()
	}
}

// Results returns the aggregated results after Flush has been called.
func (s *AggregateStage) Results() []AggregateResult {
	return s.result
}

// WriteResults writes all aggregated results to w in a human-readable form.
func (s *AggregateStage) WriteResults(w io.Writer) error {
	s.Flush()
	for _, r := range s.result {
		var line string
		if s.agg.mode == AggCount {
			line = fmt.Sprintf("%s\t%d\n", r.Key, r.Count)
		} else {
			line = fmt.Sprintf("%s\t%g\t(n=%d)\n", r.Key, r.Value, r.Count)
		}
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}
