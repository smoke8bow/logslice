package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeMergeLine(fields map[string]interface{}) parser.LogLine {
	return parser.LogLine{
		Raw:    "test line",
		Fields: fields,
	}
}

func TestNewFieldMerger_TooFewSources(t *testing.T) {
	_, err := NewFieldMerger([]string{"a"}, "out", " ", false)
	if err == nil {
		t.Fatal("expected error for fewer than two sources")
	}
}

func TestNewFieldMerger_BlankSource(t *testing.T) {
	_, err := NewFieldMerger([]string{"a", "  "}, "out", " ", false)
	if err == nil {
		t.Fatal("expected error for blank source field")
	}
}

func TestNewFieldMerger_BlankDest(t *testing.T) {
	_, err := NewFieldMerger([]string{"a", "b"}, "", " ", false)
	if err == nil {
		t.Fatal("expected error for blank destination field")
	}
}

func TestNewFieldMerger_Valid(t *testing.T) {
	fm, err := NewFieldMerger([]string{"first", "last"}, "full_name", " ", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm == nil {
		t.Fatal("expected non-nil FieldMerger")
	}
}

func TestFieldMerger_MergesFields(t *testing.T) {
	fm, _ := NewFieldMerger([]string{"first", "last"}, "full_name", " ", false)
	line := makeMergeLine(map[string]interface{}{"first": "Jane", "last": "Doe"})
	out, keep := fm.Process(line)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	if got, ok := out.Fields["full_name"]; !ok || got != "Jane Doe" {
		t.Fatalf("expected full_name=\"Jane Doe\", got %v", got)
	}
	// Source fields should still be present when remove=false.
	if _, ok := out.Fields["first"]; !ok {
		t.Error("expected source field 'first' to remain")
	}
}

func TestFieldMerger_RemovesSourceFields(t *testing.T) {
	fm, _ := NewFieldMerger([]string{"first", "last"}, "full_name", "-", true)
	line := makeMergeLine(map[string]interface{}{"first": "John", "last": "Smith"})
	out, _ := fm.Process(line)
	if _, ok := out.Fields["first"]; ok {
		t.Error("expected source field 'first' to be removed")
	}
	if _, ok := out.Fields["last"]; ok {
		t.Error("expected source field 'last' to be removed")
	}
	if got := out.Fields["full_name"]; got != "John-Smith" {
		t.Fatalf("expected full_name=\"John-Smith\", got %v", got)
	}
}

func TestFieldMerger_PlainLinePassthrough(t *testing.T) {
	fm, _ := NewFieldMerger([]string{"a", "b"}, "c", ":", false)
	line := parser.LogLine{Raw: "unstructured log line"}
	out, keep := fm.Process(line)
	if !keep {
		t.Fatal("expected plain line to be kept")
	}
	if out.Raw != "unstructured log line" {
		t.Errorf("unexpected raw change: %q", out.Raw)
	}
}

func TestFieldMerger_MissingSourcesSkipped(t *testing.T) {
	fm, _ := NewFieldMerger([]string{"x", "y"}, "z", "_", false)
	line := makeMergeLine(map[string]interface{}{"other": "value"})
	out, keep := fm.Process(line)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	// No source fields present, dest should not be set.
	if _, ok := out.Fields["z"]; ok {
		t.Error("expected 'z' not to be set when sources are missing")
	}
}

// TestFieldMerger_PartialSourcesSkipped verifies that when only a subset of
// source fields are present the merge is skipped and no destination is written.
func TestFieldMerger_PartialSourcesSkipped(t *testing.T) {
	fm, _ := NewFieldMerger([]string{"first", "last"}, "full_name", " ", false)
	line := makeMergeLine(map[string]interface{}{"first": "OnlyFirst"})
	out, keep := fm.Process(line)
	if !keep {
		t.Fatal("expected line to be kept")
	}
	// Only one of the two source fields is present; dest should not be set.
	if _, ok := out.Fields["full_name"]; ok {
		t.Error("expected 'full_name' not to be set when only partial sources are present")
	}
}
