package processor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nicholasgasior/logslice/internal/output"
	"github.com/nicholasgasior/logslice/internal/reader"
)

func TestNewJSONPathStage_InvalidPath(t *testing.T) {
	_, err := NewJSONPathStage("", "out")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewJSONPathStage_Valid(t *testing.T) {
	s, err := NewJSONPathStage("request.host", "host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil stage")
	}
}

func TestJSONPathStage_Run_ExtractsField(t *testing.T) {
	input := `{"request.host":"example.com","status":"200"}` + "\n"
	r, err := reader.NewLineReader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("reader: %v", err)
	}

	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	s, _ := NewJSONPathStage("request.host", "host")
	if runErr := s.Run(r, w, "raw"); runErr != nil {
		t.Fatalf("Run error: %v", runErr)
	}

	out := buf.String()
	if out == "" {
		t.Fatal("expected output, got empty string")
	}
}

func TestJSONPathStage_Run_EmptyInput(t *testing.T) {
	r, _ := reader.NewLineReader(strings.NewReader(""))
	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	s, _ := NewJSONPathStage("a.b", "out")
	if err := s.Run(r, w, "raw"); err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestJSONPathStage_Run_PlainLinePassthrough(t *testing.T) {
	input := "just a plain log line\n"
	r, _ := reader.NewLineReader(strings.NewReader(input))
	var buf bytes.Buffer
	w := output.NewWriter(&buf)

	s, _ := NewJSONPathStage("request.host", "host")
	if err := s.Run(r, w, "raw"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Plain lines have no fields; they should still be written through.
	if !strings.Contains(buf.String(), "plain log line") {
		t.Errorf("expected plain line in output, got %q", buf.String())
	}
}
