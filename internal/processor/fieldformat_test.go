package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeFormatLine(fields map[string]string) parser.LogLine {
	return parser.LogLine{
		Raw:    "raw log line",
		Fields: fields,
	}
}

func TestNewFieldFormatter_EmptyField(t *testing.T) {
	_, err := NewFieldFormatter("", "%.2f")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewFieldFormatter_EmptyFormat(t *testing.T) {
	_, err := NewFieldFormatter("latency", "")
	if err == nil {
		t.Fatal("expected error for empty format string")
	}
}

func TestNewFieldFormatter_ValidFormat(t *testing.T) {
	ff, err := NewFieldFormatter("latency", "%.2f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ff == nil {
		t.Fatal("expected non-nil FieldFormatter")
	}
}

func TestFieldFormatter_FloatField(t *testing.T) {
	ff, _ := NewFieldFormatter("latency", "%.3f")
	line := makeFormatLine(map[string]string{"latency": "1.5", "level": "info"})
	out := ff.Process(line)
	if got := out.Fields["latency"]; got != "1.500" {
		t.Errorf("expected '1.500', got %q", got)
	}
	if out.Fields["level"] != "info" {
		t.Error("unrelated field should be preserved")
	}
}

func TestFieldFormatter_IntegerField(t *testing.T) {
	ff, _ := NewFieldFormatter("code", "%05d")
	line := makeFormatLine(map[string]string{"code": "42"})
	out := ff.Process(line)
	if got := out.Fields["code"]; got != "00042" {
		t.Errorf("expected '00042', got %q", got)
	}
}

func TestFieldFormatter_StringFallback(t *testing.T) {
	ff, _ := NewFieldFormatter("env", "%q")
	line := makeFormatLine(map[string]string{"env": "production"})
	out := ff.Process(line)
	if got := out.Fields["env"]; got != `"production"` {
		t.Errorf("expected quoted string, got %q", got)
	}
}

func TestFieldFormatter_MissingField_Passthrough(t *testing.T) {
	ff, _ := NewFieldFormatter("missing", "%.2f")
	line := makeFormatLine(map[string]string{"level": "warn"})
	out := ff.Process(line)
	if _, ok := out.Fields["missing"]; ok {
		t.Error("missing field should not be created")
	}
}

func TestFieldFormatter_NilFields_Passthrough(t *testing.T) {
	ff, _ := NewFieldFormatter("latency", "%.2f")
	line := parser.LogLine{Raw: "plain text", Fields: nil}
	out := ff.Process(line)
	if out.Raw != "plain text" {
		t.Error("raw line should be unchanged")
	}
}

func TestFieldFormatter_DoesNotMutateOriginal(t *testing.T) {
	ff, _ := NewFieldFormatter("code", "%05d")
	orig := map[string]string{"code": "7"}
	line := makeFormatLine(orig)
	ff.Process(line)
	if orig["code"] != "7" {
		t.Error("original fields map should not be mutated")
	}
}
