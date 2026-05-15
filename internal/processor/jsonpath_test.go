package processor

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeJSONPathLine(fields map[string]string) *parser.LogLine {
	return &parser.LogLine{
		Raw:    "raw line",
		Fields: fields,
	}
}

func TestNewJSONPathExtractor_EmptyPath(t *testing.T) {
	_, err := NewJSONPathExtractor("", "out")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewJSONPathExtractor_EmptyOutput(t *testing.T) {
	_, err := NewJSONPathExtractor("a.b", "")
	if err == nil {
		t.Fatal("expected error for empty output field")
	}
}

func TestNewJSONPathExtractor_EmptySegment(t *testing.T) {
	_, err := NewJSONPathExtractor("a..b", "out")
	if err == nil {
		t.Fatal("expected error for empty path segment")
	}
}

func TestNewJSONPathExtractor_Valid(t *testing.T) {
	e, err := NewJSONPathExtractor("request.host", "host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestJSONPathExtractor_NilLine(t *testing.T) {
	e, _ := NewJSONPathExtractor("a.b", "out")
	result, err := e.Process(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for nil input")
	}
}

func TestJSONPathExtractor_NoFields(t *testing.T) {
	e, _ := NewJSONPathExtractor("a.b", "out")
	line := &parser.LogLine{Raw: "plain text"}
	result, err := e.Process(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != line {
		t.Fatal("expected same line returned when no fields")
	}
}

func TestJSONPathExtractor_FlatDottedKey(t *testing.T) {
	e, _ := NewJSONPathExtractor("request.host", "host")
	line := makeJSONPathLine(map[string]string{
		"request.host": "example.com",
		"status":       "200",
	})
	result, err := e.Process(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Fields["host"] != "example.com" {
		t.Errorf("expected host=example.com, got %q", result.Fields["host"])
	}
}

func TestJSONPathExtractor_MissingPath(t *testing.T) {
	e, _ := NewJSONPathExtractor("request.missing", "out")
	line := makeJSONPathLine(map[string]string{
		"request.host": "example.com",
	})
	result, err := e.Process(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.Fields["out"]; ok {
		t.Error("expected output field to be absent when path not found")
	}
}

func TestJSONPathExtractor_SingleSegment(t *testing.T) {
	e, _ := NewJSONPathExtractor("level", "lvl")
	line := makeJSONPathLine(map[string]string{"level": "info"})
	result, err := e.Process(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Fields["lvl"] != "info" {
		t.Errorf("expected lvl=info, got %q", result.Fields["lvl"])
	}
}
