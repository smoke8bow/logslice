package processor

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

func writeTempLogFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-worker-*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer f.Close()
	_, _ = f.WriteString(content)
	return f.Name()
}

func TestWorkerPool_SingleWorker(t *testing.T) {
	path := writeTempLogFile(t, "line one\nline two\nline three\n")

	var buf bytes.Buffer
	w := makeWriter(&buf)
	p := parser.NewParser()
	pf, _ := filter.NewPatternFilter(nil, nil)

	pool := NewWorkerPool(1, p, pf, nil, w, output.FormatRaw)
	pool.Run([]Job{{FilePath: path, Offset: 0}})

	got := buf.String()
	for _, line := range []string{"line one", "line two", "line three"} {
		if !strings.Contains(got, line) {
			t.Errorf("expected output to contain %q, got:\n%s", line, got)
		}
	}
}

func TestWorkerPool_MultipleWorkers(t *testing.T) {
	path := writeTempLogFile(t, "alpha\nbeta\ngamma\ndelta\n")

	var buf bytes.Buffer
	w := makeWriter(&buf)
	p := parser.NewParser()
	pf, _ := filter.NewPatternFilter(nil, nil)

	pool := NewWorkerPool(4, p, pf, nil, w, output.FormatRaw)
	pool.Run([]Job{{FilePath: path, Offset: 0}})

	got := buf.String()
	if got == "" {
		t.Error("expected non-empty output from worker pool")
	}
}

func TestWorkerPool_InvalidFile(t *testing.T) {
	var buf bytes.Buffer
	w := makeWriter(&buf)
	p := parser.NewParser()
	pf, _ := filter.NewPatternFilter(nil, nil)

	pool := NewWorkerPool(1, p, pf, nil, w, output.FormatRaw)
	// Should not panic on a missing file
	pool.Run([]Job{{FilePath: "/nonexistent/path/file.log", Offset: 0}})

	if buf.Len() != 0 {
		t.Errorf("expected no output for invalid file, got: %s", buf.String())
	}
}

func TestNewWorkerPool_MinWorkers(t *testing.T) {
	p := parser.NewParser()
	pf, _ := filter.NewPatternFilter(nil, nil)
	var buf bytes.Buffer
	w := makeWriter(&buf)

	pool := NewWorkerPool(0, p, pf, nil, w, output.FormatRaw)
	if pool.numWorkers != 1 {
		t.Errorf("expected numWorkers=1 for zero input, got %d", pool.numWorkers)
	}
}
