package output

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriter_WriteLine(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteLine("hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	if got := buf.String(); got != "hello world\n" {
		t.Errorf("expected %q, got %q", "hello world\n", got)
	}
	if w.LineCount() != 1 {
		t.Errorf("expected line count 1, got %d", w.LineCount())
	}
}

func TestWriter_WriteBytes(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteBytes([]byte("raw bytes")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = w.Flush()
	if !strings.Contains(buf.String(), "raw bytes") {
		t.Errorf("expected output to contain 'raw bytes', got %q", buf.String())
	}
	if w.LineCount() != 1 {
		t.Errorf("expected line count 1, got %d", w.LineCount())
	}
}

func TestWriter_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	for i := 0; i < 5; i++ {
		_ = w.WriteLine("line")
	}
	_ = w.Flush()
	if w.LineCount() != 5 {
		t.Errorf("expected 5 lines, got %d", w.LineCount())
	}
}

func TestNewFileWriter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	fw, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = fw.WriteLine("file line")
	if err := fw.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if string(data) != "file line\n" {
		t.Errorf("expected %q, got %q", "file line\n", string(data))
	}
}

func TestNewFileWriter_InvalidPath(t *testing.T) {
	_, err := NewFileWriter("/nonexistent/dir/out.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
