package processor

import (
	"fmt"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

// FieldCastStage wraps FieldCaster for use in a processing pipeline.
type FieldCastStage struct {
	caster *FieldCaster
}

// NewFieldCastStage creates a pipeline stage that casts a field to the given type.
func NewFieldCastStage(field string, castType CastType) (*FieldCastStage, error) {
	c, err := NewFieldCaster(field, castType)
	if err != nil {
		return nil, err
	}
	return &FieldCastStage{caster: c}, nil
}

// Run reads all lines from lr, applies the cast, and writes results to w.
func (s *FieldCastStage) Run(lr *reader.LineReader, p *parser.Parser, w *output.Writer) error {
	for {
		raw, err := lr.ReadLine()
		if err != nil {
			break
		}
		if len(raw) == 0 {
			continue
		}
		line := p.Parse(string(raw))
		line = s.caster.Process(line)
		formatted := output.FormatLine(line, w.Format())
		if err := w.WriteBytes([]byte(formatted)); err != nil {
			return fmt.Errorf("fieldcast: write: %w", err)
		}
	}
	return nil
}
