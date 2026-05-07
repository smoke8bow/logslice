package processor

import (
	"os"
	"strings"
	"testing"

	"github.com/user/logslice/internal/config"
	"github.com/user/logslice/internal/output"
)

func writeTempParallelLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "parallel_test_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestRunParallel_SingleWorker(t *testing.T) {
	lines := []string{"line one", "line two", "line three"}
	path := writeTempParallelLog(t, lines)
	defer os.Remove(path)

	var buf strings.Builder
	w := output.NewWriter(&buf)

	cfg := &config.Config{
		InputFile: path,
		Workers:   1,
		Format:    "raw",
	}

	result, err := RunParallel(cfg, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LinesRead < 1 {
		t.Errorf("expected at least 1 line read, got %d", result.LinesRead)
	}
	if len(result.Errors) > 0 {
		t.Errorf("unexpected worker errors: %v", result.Errors)
	}
}

func TestRunParallel_MultipleWorkers(t *testing.T) {
	lines := make([]string, 20)
	for i := range lines {
		lines[i] = "log entry"
	}
	path := writeTempParallelLog(t, lines)
	defer os.Remove(path)

	var buf strings.Builder
	w := output.NewWriter(&buf)

	cfg := &config.Config{
		InputFile: path,
		Workers:   4,
		Format:    "raw",
	}

	result, err := RunParallel(cfg, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Errors) > 0 {
		t.Errorf("unexpected worker errors: %v", result.Errors)
	}
}

func TestRunParallel_MissingFile(t *testing.T) {
	var buf strings.Builder
	w := output.NewWriter(&buf)

	cfg := &config.Config{
		InputFile: "/nonexistent/path/file.log",
		Workers:   2,
		Format:    "raw",
	}

	_, err := RunParallel(cfg, w)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRunParallel_WithPatternFilter(t *testing.T) {
	lines := []string{"error: disk full", "info: all good", "error: timeout"}
	path := writeTempParallelLog(t, lines)
	defer os.Remove(path)

	var buf strings.Builder
	w := output.NewWriter(&buf)

	cfg := &config.Config{
		InputFile: path,
		Workers:   1,
		Format:    "raw",
		Include:   "error",
	}

	result, err := RunParallel(cfg, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.LinesWritten != 2 {
		t.Errorf("expected 2 lines written, got %d", result.LinesWritten)
	}
}
