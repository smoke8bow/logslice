package main

import (
	"fmt"
	"os"

	"github.com/user/logslice/internal/config"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/processor"
	"github.com/user/logslice/internal/reader"
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

	var opts []processor.Option

	if cfg.Start != "" || cfg.End != "" {
		tr, err := filter.ParseTimeRange(cfg.Start, cfg.End)
		if err != nil {
			return fmt.Errorf("time range: %w", err)
		}
		opts = append(opts, processor.WithTimeRange(tr))
	}

	if cfg.Include != "" || cfg.Exclude != "" {
		pf, err := filter.NewPatternFilter(cfg.Include, cfg.Exclude)
		if err != nil {
			return fmt.Errorf("pattern filter: %w", err)
		}
		opts = append(opts, processor.WithPatternFilter(pf))
	}

	if cfg.RateLimit > 0 {
		rl, err := processor.NewRateLimiter(cfg.RateLimit)
		if err != nil {
			return fmt.Errorf("rate limit: %w", err)
		}
		opts = append(opts, processor.WithRateLimiter(rl))
	}

	fmt, err := output.ParseFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("output format: %w", err)
	}

	w, err := output.NewFileWriter(cfg.Output, fmt)
	if err != nil {
		return fmt.Errorf("output writer: %w", err)
	}
	defer w.Close()

	if cfg.Workers > 1 {
		return processor.RunParallel(cfg.Input, w, cfg.Workers, opts...)
	}

	r, err := reader.NewLineReader(cfg.Input)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer r.Close()

	pipe := processor.New(r, w, opts...)
	stats := processor.NewStats()
	if err := pipe.Run(stats); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, stats.Summary())
	return nil
}
