package reader

import (
	"strings"
	"testing"
)

func TestScanner_EmptyInput(t *testing.T) {
	s := NewScanner(strings.NewReader(""))
	if s.Scan() {
		t.Fatal("expected no lines from empty input")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScanner_SingleLine(t *testing.T) {
	s := NewScanner(strings.NewReader("hello world"))
	if !s.Scan() {
		t.Fatal("expected a line")
	}
	line := s.Line()
	if line.Raw != "hello world" {
		t.Errorf("expected 'hello world', got %q", line.Raw)
	}
	if line.LineNum != 1 {
		t.Errorf("expected LineNum=1, got %d", line.LineNum)
	}
	if line.Offset != 0 {
		t.Errorf("expected Offset=0, got %d", line.Offset)
	}
}

func TestScanner_MultipleLines(t *testing.T) {
	input := "line one\nline two\nline three"
	s := NewScanner(strings.NewReader(input))

	expected := []string{"line one", "line two", "line three"}
	for i, want := range expected {
		if !s.Scan() {
			t.Fatalf("expected line %d", i+1)
		}
		line := s.Line()
		if line.Raw != want {
			t.Errorf("line %d: expected %q, got %q", i+1, want, line.Raw)
		}
		if line.LineNum != i+1 {
			t.Errorf("line %d: expected LineNum=%d, got %d", i+1, i+1, line.LineNum)
		}
	}
	if s.Scan() {
		t.Fatal("expected no more lines")
	}
}

func TestScanner_OffsetTracking(t *testing.T) {
	input := "abc\ndefg\nhi"
	s := NewScanner(strings.NewReader(input))

	expectedOffsets := []int64{0, 4, 9}
	for i, wantOffset := range expectedOffsets {
		if !s.Scan() {
			t.Fatalf("expected line %d", i+1)
		}
		line := s.Line()
		if line.Offset != wantOffset {
			t.Errorf("line %d: expected offset %d, got %d", i+1, wantOffset, line.Offset)
		}
	}
}

func TestScanner_NoErrorOnCleanRead(t *testing.T) {
	s := NewScanner(strings.NewReader("a\nb\nc"))
	for s.Scan() {
		s.Line()
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
