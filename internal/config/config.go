package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Config holds all parsed CLI configuration for logslice.
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
	SampleRate  float64
	SampleNth   int
	SampleMode  string
	MaxBytes    int
	Dedup       bool
	Merge       bool
	Field       string
	FieldValue  string
	FieldRegexp string
	MaskFields  string
	Transform   string
	TransformField string
	RateLimit   int
}

// Parse parses the command-line arguments and returns a Config or an error.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	output   := fs.String("output", "", "output file (default: stdout)")
	format   := fs.String("format", "raw", "output format: raw, json")
	start    := fs.String("start", "", "start timestamp (RFC3339)")
	end      := fs.String("end", "", "end timestamp (RFC3339)")
	include  := fs.String("include", "", "include pattern (regex)")
	exclude  := fs.String("exclude", "", "exclude pattern (regex)")
	level    := fs.String("level", "", "minimum log level")
	workers  := fs.Int("workers", 1, "number of parallel workers")
	sampleRate := fs.Float64("sample-rate", 0, "probabilistic sample rate (0-1)")
	sampleNth  := fs.Int("sample-nth", 0, "keep every Nth line")
	sampleMode := fs.String("sample-mode", "", "sampling mode: rate or nth")
	maxBytes   := fs.Int("max-bytes", 0, "truncate lines to max bytes")
	dedup      := fs.Bool("dedup", false, "deduplicate identical lines")
	merge      := fs.Bool("merge", false, "merge and sort multiple input files by timestamp")
	field      := fs.String("field", "", "field name for field filter")
	fieldValue := fs.String("field-value", "", "exact field value")
	fieldRegexp := fs.String("field-regexp", "", "field value regexp")
	maskFields  := fs.String("mask-fields", "", "comma-separated fields to mask")
	transform      := fs.String("transform", "", "transform to apply: upper, lower, trim, urlencode")
	transformField := fs.String("transform-field", "", "field to apply transform to")
	rateLimit := fs.Int("rate-limit", 0, "max lines per second (0 = unlimited)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	inputs := fs.Args()
	if len(inputs) == 0 {
		return nil, errors.New("at least one input file is required")
	}

	if *format != "raw" && *format != "json" {
		return nil, fmt.Errorf("invalid format %q: must be raw or json", *format)
	}

	if *workers < 1 {
		return nil, errors.New("workers must be at least 1")
	}

	if *transform != "" && *transformField == "" {
		return nil, errors.New("--transform-field is required when --transform is set")
	}

	knownTransforms := map[string]bool{"upper": true, "lower": true, "trim": true, "urlencode": true}
	if *transform != "" && !knownTransforms[strings.ToLower(*transform)] {
		return nil, fmt.Errorf("unknown transform %q", *transform)
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
		SampleRate:     *sampleRate,
		SampleNth:      *sampleNth,
		SampleMode:     *sampleMode,
		MaxBytes:       *maxBytes,
		Dedup:          *dedup,
		Merge:          *merge,
		Field:          *field,
		FieldValue:     *fieldValue,
		FieldRegexp:    *fieldRegexp,
		MaskFields:     *maskFields,
		Transform:      *transform,
		TransformField: *transformField,
		RateLimit:      *rateLimit,
	}, nil
}
