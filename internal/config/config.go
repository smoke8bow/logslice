package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Config holds the parsed CLI configuration for logslice.
type Config struct {
	InputFiles  []string
	OutputFile  string
	Format      string
	Start       string
	End         string
	Include     string
	Exclude     string
	Workers     int
	Merge       bool
}

var validFormats = map[string]bool{
	"raw":  true,
	"json": true,
}

// Parse reads command-line arguments and returns a Config or an error.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	output := fs.String("output", "", "output file (default: stdout)")
	format := fs.String("format", "raw", "output format: raw or json")
	start := fs.String("start", "", "start timestamp filter (RFC3339)")
	end := fs.String("end", "", "end timestamp filter (RFC3339)")
	include := fs.String("include", "", "include pattern (regex)")
	exclude := fs.String("exclude", "", "exclude pattern (regex)")
	workers := fs.Int("workers", 1, "number of parallel workers (>=1)")
	merge := fs.Bool("merge", false, "merge and sort multiple input files by timestamp")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	inputs := fs.Args()
	if len(inputs) == 0 {
		return nil, errors.New("at least one input file is required")
	}

	if !validFormats[*format] {
		return nil, fmt.Errorf("invalid format %q: must be raw or json", *format)
	}

	if *workers < 1 {
		return nil, errors.New("workers must be >= 1")
	}

	return &Config{
		InputFiles: inputs,
		OutputFile: *output,
		Format:     *format,
		Start:      *start,
		End:        *end,
		Include:    *include,
		Exclude:    *exclude,
		Workers:    *workers,
		Merge:      *merge,
	}, nil
}
