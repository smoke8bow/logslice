package processor

import (
	"fmt"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
	"github.com/user/logslice/internal/output"
)

// ContextConfig holds configuration for context-aware log extraction.
type ContextConfig struct {
	InputFile   string
	Pattern     string
	Before      int
	After       int
	OutputFile  string
	Format      string
}

// RunWithContext runs the pipeline with before/after context line support.
func RunWithContext(cfg ContextConfig) error {
	pf, err := filter.NewPatternFilter([]string{cfg.Pattern}, nil)
	if err != nil {
		return fmt.Errorf("contextlines: invalid pattern: %w", err)
	}

	cl, err := NewContextLines(cfg.Before, cfg.After)
	if err != nil {
		return fmt.Errorf("contextlines: %w", err)
	}

	fmt, err := output.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("contextlines: invalid format: %w", err)
	}

	var w *output.Writer
	if cfg.OutputFile != "" {
		w, err = output.NewFileWriter(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("contextlines: open output: %w", err)
		}
		defer w.Close()
	} else {
		w = output.NewWriter(nil)
	}

	lr, err := reader.NewLineReader(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("contextlines: open input: %w", err)
	}
	defer lr.Close()

	p := parser.NewParser()

	for {
		line, err := lr.ReadLine()
		if err != nil {
			break
		}
		parsed := p.Parse(line)
		matched := pf.Match(parsed)
		emitted := cl.Process(parsed, matched)
		for _, el := range emitted {
			if werr := w.WriteLine(output.FormatLine(el, fmt)); werr != nil {
				return werr
			}
		}
	}
	return nil
}
