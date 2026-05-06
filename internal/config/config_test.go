package config

import (
	"testing"
)

func TestParse_MissingInputFile(t *testing.T) {
	_, err := Parse([]string{})
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestParse_MinimalValid(t *testing.T) {
	cfg, err := Parse([]string{"-f", "app.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "app.log" {
		t.Errorf("expected InputFile=app.log, got %s", cfg.InputFile)
	}
	if cfg.Format != "raw" {
		t.Errorf("expected default format=raw, got %s", cfg.Format)
	}
	if cfg.Workers != 1 {
		t.Errorf("expected default workers=1, got %d", cfg.Workers)
	}
}

func TestParse_AllFlags(t *testing.T) {
	cfg, err := Parse([]string{
		"-f", "input.log",
		"-o", "out.log",
		"-start", "2024-01-01T00:00:00Z",
		"-end", "2024-01-02T00:00:00Z",
		"-format", "json",
		"-workers", "4",
		"-include", "ERROR",
		"-include", "WARN",
		"-exclude", "DEBUG",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFile != "out.log" {
		t.Errorf("expected out.log, got %s", cfg.OutputFile)
	}
	if cfg.Workers != 4 {
		t.Errorf("expected 4 workers, got %d", cfg.Workers)
	}
	if len(cfg.Include) != 2 {
		t.Errorf("expected 2 include patterns, got %d", len(cfg.Include))
	}
	if len(cfg.Exclude) != 1 {
		t.Errorf("expected 1 exclude pattern, got %d", len(cfg.Exclude))
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]string{"-f", "app.log", "-format", "xml"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestParse_InvalidWorkers(t *testing.T) {
	_, err := Parse([]string{"-f", "app.log", "-workers", "0"})
	if err == nil {
		t.Fatal("expected error for workers=0")
	}
}

func TestParse_InvalidChunkSize(t *testing.T) {
	_, err := Parse([]string{"-f", "app.log", "-chunk-size", "0"})
	if err == nil {
		t.Fatal("expected error for chunk-size=0")
	}
}
