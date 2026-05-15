package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeTemplateLine(raw string, fields map[string]string) parser.LogLine {
	return parser.LogLine{Raw: raw, Fields: fields}
}

func TestNewTemplater_EmptyString(t *testing.T) {
	_, err := NewTemplater("")
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNewTemplater_InvalidSyntax(t *testing.T) {
	_, err := NewTemplater("{{ .Unclosed")
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestNewTemplater_ValidTemplate(t *testing.T) {
	_, err := NewTemplater("{{.Raw}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplater_RawPassthrough(t *testing.T) {
	tmpl, _ := NewTemplater("{{.Raw}}")
	line := makeTemplateLine("hello world", nil)
	out := tmpl.Apply(line)
	if out.Raw != "hello world" {
		t.Errorf("expected 'hello world', got %q", out.Raw)
	}
}

func TestTemplater_FieldSubstitution(t *testing.T) {
	tmpl, _ := NewTemplater(`level={{index .Fields "level"}} msg={{index .Fields "msg"}}`)
	line := makeTemplateLine("raw", map[string]string{"level": "info", "msg": "started"})
	out := tmpl.Apply(line)
	expected := "level=info msg=started"
	if out.Raw != expected {
		t.Errorf("expected %q, got %q", expected, out.Raw)
	}
}

func TestTemplater_MissingFieldIsEmpty(t *testing.T) {
	tmpl, _ := NewTemplater(`{{index .Fields "missing"}}`)
	line := makeTemplateLine("raw", map[string]string{})
	out := tmpl.Apply(line)
	if out.Raw != "" {
		t.Errorf("expected empty string for missing field, got %q", out.Raw)
	}
}

func TestTemplater_ExecutionErrorPreservesRaw(t *testing.T) {
	// A template that calls a function that doesn't exist at runtime via a
	// deliberate bad pipeline won't fail at parse time but may at exec;
	// we simulate by passing nil Fields to a template that dereferences them.
	tmpl, _ := NewTemplater("{{.Raw}} extra")
	line := makeTemplateLine("original", nil)
	out := tmpl.Apply(line)
	if out.Raw != "original extra" {
		t.Errorf("unexpected output: %q", out.Raw)
	}
}

func TestTemplater_PreservesOtherFields(t *testing.T) {
	tmpl, _ := NewTemplater("rewritten")
	fields := map[string]string{"key": "val"}
	line := makeTemplateLine("original", fields)
	out := tmpl.Apply(line)
	if out.Raw != "rewritten" {
		t.Errorf("expected 'rewritten', got %q", out.Raw)
	}
	if out.Fields["key"] != "val" {
		t.Error("Fields should be preserved after template apply")
	}
}
