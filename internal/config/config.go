package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Config holds all CLI-parsed configuration for a logslice run.
type Config struct {
	InputFile  string
	OutputFile string
	Start      string
	End        string
	Include    []string
	Exclude    []string
	Format     string
	Workers    int
	ChunkSize  int64
}

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// Parse parses command-line arguments into a Config.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	cfg := &Config{}
	var include, exclude multiFlag

	fs.StringVar(&cfg.InputFile, "f", "", "input log file (required)")
	fs.StringVar(&cfg.OutputFile, "o", "", "output file (default: stdout)")
	fs.StringVar(&cfg.Start, "start", "", "start timestamp (RFC3339 or common log formats)")
	fs.StringVar(&cfg.End, "end", "", "end timestamp")
	fs.StringVar(&cfg.Format, "format", "raw", "output format: raw|json")
	fs.IntVar(&cfg.Workers, "workers", 1, "number of parallel workers")
	fs.Int64Var(&cfg.ChunkSize, "chunk-size", 64*1024*1024, "chunk size in bytes for parallel reads")
	fs.Var(&include, "include", "include pattern (regex); may be repeated")
	fs.Var(&exclude, "exclude", "exclude pattern (regex); may be repeated")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg.Include = []string(include)
	cfg.Exclude = []string(exclude)

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.InputFile == "" {
		return errors.New("input file (-f) is required")
	}
	if c.Workers < 1 {
		return fmt.Errorf("workers must be >= 1, got %d", c.Workers)
	}
	if c.ChunkSize < 1 {
		return fmt.Errorf("chunk-size must be >= 1, got %d", c.ChunkSize)
	}
	valid := map[string]bool{"raw": true, "json": true}
	if !valid[c.Format] {
		return fmt.Errorf("unknown format %q; must be raw or json", c.Format)
	}
	return nil
}
