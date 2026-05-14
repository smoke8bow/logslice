package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds all parsed CLI configuration for a logslice run.
type Config struct {
	InputFiles  []string
	OutputFile  string
	Format      string
	Start       string
	End         string
	Include     string
	Exclude     string
	Level       string
	Workers     int
	SampleNth   int
	Dedup       bool
	Highlight   string // comma-separated patterns
	HighlightColor string
}

// Parse reads os.Args via the provided FlagSet and returns a validated Config.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	output := fs.String("output", "", "write output to file (default: stdout)")
	format := fs.String("format", "raw", "output format: raw or json")
	start := fs.String("start", "", "start timestamp (RFC3339)")
	end := fs.String("end", "", "end timestamp (RFC3339)")
	include := fs.String("include", "", "include lines matching regex")
	exclude := fs.String("exclude", "", "exclude lines matching regex")
	level := fs.String("level", "", "minimum log level (debug|info|warn|error)")
	workers := fs.Int("workers", 1, "number of parallel workers")
	sampleNth := fs.Int("sample-nth", 0, "keep every Nth line (0 = disabled)")
	dedup := fs.Bool("dedup", false, "deduplicate identical lines")
	highlight := fs.String("highlight", "", "comma-separated regex patterns to highlight")
	highlightColor := fs.String("highlight-color", "cyan", "highlight color: red, yellow, cyan, green")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	inputs := fs.Args()
	if len(inputs) == 0 {
		return nil, errors.New("at least one input file is required")
	}

	*format = strings.ToLower(strings.TrimSpace(*format))
	if *format != "raw" && *format != "json" {
		return nil, fmt.Errorf("invalid format %q: must be raw or json", *format)
	}

	if *workers < 1 {
		return nil, fmt.Errorf("workers must be >= 1, got %d", *workers)
	}

	for _, f := range inputs {
		if _, err := os.Stat(f); err != nil {
			return nil, fmt.Errorf("input file %q: %w", f, err)
		}
	}

	return &Config{
		InputFiles:     inputs,
		OutputFile:     *output,
		Format:         *format,
		Start:          *start,
		End:            *end,
		Include:        *include,
		Exclude:        *exclude,
		Level:          *level,
		Workers:        *workers,
		SampleNth:      *sampleNth,
		Dedup:          *dedup,
		Highlight:      *highlight,
		HighlightColor: *highlightColor,
	}, nil
}
