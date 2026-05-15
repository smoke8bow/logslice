package processor

import (
	"fmt"

	"github.com/nicholasgasior/logslice/internal/parser"
	"github.com/nicholasgasior/logslice/internal/reader"
	"github.com/nicholasgasior/logslice/internal/output"
)

// JSONPathStage wraps a JSONPathExtractor for use in a processing pipeline.
type JSONPathStage struct {
	extractor *JSONPathExtractor
}

// NewJSONPathStage constructs a pipeline stage that extracts a nested field.
func NewJSONPathStage(dotPath, outputField string) (*JSONPathStage, error) {
	e, err := NewJSONPathExtractor(dotPath, outputField)
	if err != nil {
		return nil, err
	}
	return &JSONPathStage{extractor: e}, nil
}

// Run processes all lines from the reader through the extractor and writes
// matching output to the writer.
func (s *JSONPathStage) Run(r *reader.LineReader, w *output.Writer, fmt_ string) error {
	p := parser.NewParser()
	for {
		raw, err := r.ReadLine()
		if err != nil {
			break
		}
		if raw == "" {
			continue
		}
		line := p.Parse(raw)
		result, procErr := s.extractor.Process(line)
		if procErr != nil {
			return fmt.Errorf("jsonpath stage: %w", procErr)
		}
		if result == nil {
			continue
		}
		formatted := output.FormatLine(result, fmt_)
		if writeErr := w.WriteLine(formatted); writeErr != nil {
			return fmt.Errorf("jsonpath stage write: %w", writeErr)
		}
	}
	return nil
}
