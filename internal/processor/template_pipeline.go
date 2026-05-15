package processor

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// TemplateStage wraps a Templater so it satisfies a simple stage interface
// compatible with the existing pipeline: it receives a LogLine, applies the
// template, and returns the modified line.
type TemplateStage struct {
	templater *Templater
}

// NewTemplateStage creates a TemplateStage from a raw template string.
func NewTemplateStage(tmplStr string) (*TemplateStage, error) {
	tmpl, err := NewTemplater(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template stage: %w", err)
	}
	return &TemplateStage{templater: tmpl}, nil
}

// Process applies the template to the line and always passes it through
// (template stages never drop lines, only rewrite them).
func (ts *TemplateStage) Process(line parser.LogLine) (parser.LogLine, bool) {
	return ts.templater.Apply(line), true
}

// String returns a human-readable description of the stage.
func (ts *TemplateStage) String() string {
	return fmt.Sprintf("TemplateStage(tmpl=%q)", ts.templater.tmpl.Root.String())
}
