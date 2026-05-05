package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempLog: %v", err)
	}
	return p
}

func TestNewLineReader_NotFound(t *testing.T) {
	_, err := NewLineReader("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLineReader_ReadAllLines(t *testing.T) {
	path := writeTempLog(t, "line1\nline2\nline3\n")
	r, err := NewLineReader(path)
	if err != nil {
		t.Fatalf("NewLineReader: %v", err)
	}
	defer r.Close()

	var lines []string
	for r.Scan() {
		lines = append(lines, r.Text())
	}
	if r.Err() != nil {
		t.Fatalf("scan error: %v", r.Err())
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[1] != "line2" {
		t.Errorf("expected 'line2', got %q", lines[1])
	}
}

func TestLineReader_EmptyFile(t *testing.T) {
	path := writeTempLog(t, "")
	r, err := NewLineReader(path)
	if err != nil {
		t.Fatalf("NewLineReader: %v", err)
	}
	defer r.Close()

	count := 0
	for r.Scan() {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 lines from empty file, got %d", count)
	}
}

func TestNewLineReaderAt_Offset(t *testing.T) {
	// "line1\n" is 6 bytes; offset 6 should start at "line2"
	path := writeTempLog(t, "line1\nline2\nline3\n")
	r, err := NewLineReaderAt(path, 6)
	if err != nil {
		t.Fatalf("NewLineReaderAt: %v", err)
	}
	defer r.Close()

	if !r.Scan() {
		t.Fatal("expected at least one line")
	}
	if got := r.Text(); got != "line2" {
		t.Errorf("expected 'line2', got %q", got)
	}
}

func TestLineReader_Path(t *testing.T) {
	path := writeTempLog(t, "hello\n")
	r, err := NewLineReader(path)
	if err != nil {
		t.Fatalf("NewLineReader: %v", err)
	}
	defer r.Close()
	if r.Path() != path {
		t.Errorf("expected path %q, got %q", path, r.Path())
	}
}
