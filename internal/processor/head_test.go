package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeHeadLine(raw string) *parser.LogLine {
	return &parser.LogLine{Raw: raw}
}

// --- HeadFilter tests ---

func TestNewHeadFilter_InvalidN(t *testing.T) {
	_, err := NewHeadFilter(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestHeadFilter_AllowsFirstN(t *testing.T) {
	h, err := NewHeadFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 3; i++ {
		if !h.Apply(makeHeadLine("line")) {
			t.Fatalf("expected line %d to be kept", i+1)
		}
	}
	if h.Apply(makeHeadLine("line")) {
		t.Fatal("expected 4th line to be dropped")
	}
}

func TestHeadFilter_Reset(t *testing.T) {
	h, _ := NewHeadFilter(2)
	h.Apply(makeHeadLine("a"))
	h.Apply(makeHeadLine("b"))
	h.Reset()
	if !h.Apply(makeHeadLine("c")) {
		t.Fatal("expected line to be kept after reset")
	}
}

// --- TailFilter tests ---

func TestNewTailFilter_InvalidN(t *testing.T) {
	_, err := NewTailFilter(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestTailFilter_ReturnsLastN(t *testing.T) {
	tf, err := NewTailFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 5; i++ {
		tf.Collect(makeHeadLine(string(rune('a' + i))))
	}
	lines := tf.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0].Raw != "c" || lines[1].Raw != "d" || lines[2].Raw != "e" {
		t.Fatalf("unexpected tail lines: %v", lines)
	}
}

func TestTailFilter_Reset(t *testing.T) {
	tf, _ := NewTailFilter(3)
	tf.Collect(makeHeadLine("x"))
	tf.Reset()
	if len(tf.Lines()) != 0 {
		t.Fatal("expected empty buffer after reset")
	}
}

func TestTailFilter_FewerThanN(t *testing.T) {
	tf, _ := NewTailFilter(10)
	tf.Collect(makeHeadLine("only"))
	lines := tf.Lines()
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}
