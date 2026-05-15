package processor

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeHighlightLine(raw string) parser.LogLine {
	return parser.LogLine{Raw: raw}
}

func TestNewHighlighter_NoPatterns(t *testing.T) {
	_, err := NewHighlighter([]string{}, "red")
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewHighlighter_InvalidPattern(t *testing.T) {
	_, err := NewHighlighter([]string{"[invalid"}, "red")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNewHighlighter_InvalidColor(t *testing.T) {
	_, err := NewHighlighter([]string{"foo"}, "purple")
	if err == nil {
		t.Fatal("expected error for unknown color")
	}
}

func TestNewHighlighter_DefaultColor(t *testing.T) {
	h, err := NewHighlighter([]string{"foo"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.color != colorCyan {
		t.Errorf("expected cyan default, got %q", h.color)
	}
}

func TestHighlighter_ApplyColorizes(t *testing.T) {
	h, err := NewHighlighter([]string{"error"}, "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := makeHighlightLine("something error occurred")
	out := h.Apply(line)

	if !strings.Contains(out.Raw, colorRed+"error"+colorReset) {
		t.Errorf("expected ANSI-wrapped match, got: %q", out.Raw)
	}
}

func TestHighlighter_NoMatchUnchanged(t *testing.T) {
	h, err := NewHighlighter([]string{"panic"}, "yellow")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw := "everything is fine"
	out := h.Apply(makeHighlightLine(raw))
	if out.Raw != raw {
		t.Errorf("expected unchanged line, got: %q", out.Raw)
	}
}

func TestHighlighter_MultiplePatterns(t *testing.T) {
	h, err := NewHighlighter([]string{"foo", "bar"}, "green")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := h.Apply(makeHighlightLine("foo and bar are here"))
	if !strings.Contains(out.Raw, colorGreen+"foo"+colorReset) {
		t.Errorf("expected foo highlighted, got: %q", out.Raw)
	}
	if !strings.Contains(out.Raw, colorGreen+"bar"+colorReset) {
		t.Errorf("expected bar highlighted, got: %q", out.Raw)
	}
}

func TestHighlighter_PreservesOtherFields(t *testing.T) {
	h, _ := NewHighlighter([]string{"x"}, "cyan")
	line := parser.LogLine{
		Raw:    "x marks the spot",
		Fields: map[string]string{"key": "val"},
	}
	out := h.Apply(line)
	if out.Fields["key"] != "val" {
		t.Errorf("expected fields preserved, got: %v", out.Fields)
	}
}

func TestHighlighter_EmptyLine(t *testing.T) {
	h, err := NewHighlighter([]string{"error"}, "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := h.Apply(makeHighlightLine(""))
	if out.Raw != "" {
		t.Errorf("expected empty line to remain empty, got: %q", out.Raw)
	}
}
