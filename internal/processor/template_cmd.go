package processor

import (
	"fmt"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

// TemplateConfig holds parameters for RunTemplate.
type TemplateConfig struct {
	InputFile  string
	OutputFile string
	Template   string
	Format     string
}

// RunTemplate reads lines from InputFile, applies the Go template to each
// line, and writes the result to OutputFile (or stdout if empty).
func RunTemplate(cfg TemplateConfig) error {
	tmpl, err := NewTemplater(cfg.Template)
	if err != nil {
		return fmt.Errorf("template cmd: %w", err)
	}

	lr, err := reader.NewLineReader(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("template cmd: open input: %w", err)
	}
	defer lr.Close()

	var w *output.Writer
	if cfg.OutputFile == "" {
		w = output.NewWriter(nil)
	} else {
		w, err = output.NewFileWriter(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("template cmd: open output: %w", err)
		}
		defer w.Close()
	}

	fmt, err := output.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("template cmd: %w", err)
	}

	p := parser.NewParser()
	for {
		line, ok := lr.Next()
		if !ok {
			break
		}
		parsed := p.Parse(line)
		parsed = tmpl.Apply(parsed)
		w.WriteLine(output.FormatLine(parsed, fmt))
	}
	return lr.Err()
}
