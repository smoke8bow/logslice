package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeAddLine(fields map[string]string) *parser.LogLine {
	return &parser.LogLine{
		Raw:    "level=info msg=boot",
		Fields: fields,
	}
}

func TestNewFieldAdder_EmptyPairs(t *testing.T) {
	_, err := NewFieldAdder([]string{}, false)
	if err == nil {
		t.Fatal("expected error for empty pairs")
	}
}

func TestNewFieldAdder_InvalidFormat(t *testing.T) {
	_, err := NewFieldAdder([]string{"noequalssign"}, false)
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewFieldAdder_BlankKey(t *testing.T) {
	_, err := NewFieldAdder([]string{"=value"}, false)
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestFieldAdder_AddsField(t *testing.T) {
	a, err := NewFieldAdder([]string{"env=production"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeAddLine(map[string]string{"level": "info"})
	out := a.Process(line)
	if out.Fields["env"] != "production" {
		t.Errorf("expected env='production', got %q", out.Fields["env"])
	}
}

func TestFieldAdder_NoOverwriteByDefault(t *testing.T) {
	a, err := NewFieldAdder([]string{"level=debug"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeAddLine(map[string]string{"level": "info"})
	out := a.Process(line)
	if out.Fields["level"] != "info" {
		t.Errorf("expected level to remain 'info', got %q", out.Fields["level"])
	}
}

func TestFieldAdder_OverwriteWhenEnabled(t *testing.T) {
	a, err := NewFieldAdder([]string{"level=debug"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeAddLine(map[string]string{"level": "info"})
	out := a.Process(line)
	if out.Fields["level"] != "debug" {
		t.Errorf("expected level='debug', got %q", out.Fields["level"])
	}
}

func TestFieldAdder_InitialisesNilFields(t *testing.T) {
	a, err := NewFieldAdder([]string{"env=staging"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := &parser.LogLine{Raw: "plain text", Fields: nil}
	out := a.Process(line)
	if out.Fields["env"] != "staging" {
		t.Errorf("expected env='staging', got %q", out.Fields["env"])
	}
	if out.Raw != "plain text" {
		t.Error("raw value should be preserved")
	}
}

func TestFieldAdder_NilLine(t *testing.T) {
	a, _ := NewFieldAdder([]string{"env=prod"}, false)
	out := a.Process(nil)
	if out != nil {
		t.Error("expected nil for nil input")
	}
}
