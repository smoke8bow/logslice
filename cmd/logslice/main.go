package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/processor"
	"github.com/yourorg/logslice/internal/reader"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := config.Parse(args)
	if err != nil {
		return err
	}

	// Build time range filter (optional).
	var timeRange *filter.TimeRange
	if cfg.Start != "" || cfg.End != "" {
		tr, err := filter.ParseTimeRange(cfg.Start, cfg.End)
		if err != nil {
			return fmt.Errorf("time range: %w", err)
		}
		timeRange = tr
	}

	// Build pattern filter (optional).
	var patFilter *filter.PatternFilter
	if len(cfg.Include) > 0 || len(cfg.Exclude) > 0 {
		pf, err := filter.NewPatternFilter(cfg.Include, cfg.Exclude)
		if err != nil {
			return fmt.Errorf("pattern filter: %w", err)
		}
		patFilter = pf
	}

	// Open input.
	lr, err := reader.NewLineReader(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer lr.Close()

	// Open output.
	var w *output.Writer
	if cfg.OutputFile != "" {
		w, err = output.NewFileWriter(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("open output: %w", err)
		}
		defer w.Close()
	} else {
		w = output.NewWriter(os.Stdout)
	}

	fmt, err := output.ParseFormat(cfg.Format)
	if err != nil {
		return err
	}

	// Run pipeline.
	p := processor.New(lr, w, fmt, timeRange, patFilter)
	return p.Run()
}
