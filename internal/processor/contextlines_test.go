package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeCtxLine(raw string) parser.LogLine {
	return parser.LogLine{Raw: raw}
}

func TestNewContextLines_InvalidBefore(t *testing.T) {
	_, err := NewContextLines(-1, 0)
	if err == nil {
		t.Fatal("expected error for negative before")
	}
}

func TestNewContextLines_InvalidAfter(t *testing.T) {
	_, err := NewContextLines(0, -1)
	if err == nil {
		t.Fatal("expected error for negative after")
	}
}

func TestContextLines_NoContext(t *testing.T) {
	cl, _ := NewContextLines(0, 0)
	lines := []parser.LogLine{makeCtxLine("a"), makeCtxLine("b"), makeCtxLine("c")}
	matches := []bool{false, true, false}
	var got []parser.LogLine
	for i, l := range lines {
		got = append(got, cl.Process(l, matches[i])...)
	}
	if len(got) != 1 || got[0].Raw != "b" {
		t.Fatalf("expected only matched line, got %v", got)
	}
}

func TestContextLines_BeforeContext(t *testing.T) {
	cl, _ := NewContextLines(2, 0)
	lines := []parser.LogLine{
		makeCtxLine("1"), makeCtxLine("2"), makeCtxLine("3"), makeCtxLine("4"),
	}
	matches := []bool{false, false, true, false}
	var got []parser.LogLine
	for i, l := range lines {
		got = append(got, cl.Process(l, matches[i])...)
	}
	raws := make([]string, len(got))
	for i, g := range got {
		raws[i] = g.Raw
	}
	expected := []string{"2", "3"}
	if len(raws) != len(expected) {
		t.Fatalf("expected %v got %v", expected, raws)
	}
	for i := range expected {
		if raws[i] != expected[i] {
			t.Errorf("pos %d: want %s got %s", i, expected[i], raws[i])
		}
	}
}

func TestContextLines_AfterContext(t *testing.T) {
	cl, _ := NewContextLines(0, 2)
	lines := []parser.LogLine{
		makeCtxLine("a"), makeCtxLine("b"), makeCtxLine("c"), makeCtxLine("d"),
	}
	matches := []bool{false, true, false, false}
	var got []parser.LogLine
	for i, l := range lines {
		got = append(got, cl.Process(l, matches[i])...)
	}
	expected := []string{"b", "c", "d"}
	if len(got) != len(expected) {
		t.Fatalf("expected %v, got %d lines", expected, len(got))
	}
	for i, e := range expected {
		if got[i].Raw != e {
			t.Errorf("pos %d: want %s got %s", i, e, got[i].Raw)
		}
	}
}

func TestContextLines_Reset(t *testing.T) {
	cl, _ := NewContextLines(2, 2)
	cl.Process(makeCtxLine("x"), false)
	cl.Process(makeCtxLine("y"), false)
	cl.Reset()
	// After reset, before-buffer should be empty; a match should yield only itself
	out := cl.Process(makeCtxLine("match"), true)
	if len(out) != 1 || out[0].Raw != "match" {
		t.Fatalf("expected only match after reset, got %v", out)
	}
}
