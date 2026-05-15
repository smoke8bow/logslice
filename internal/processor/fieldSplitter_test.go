package processor

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeSplitLine(fields map[string]interface{}) parser.LogLine {
	return parser.LogLine{
		Raw:    "raw line",
		Fields: fields,
	}
}

func TestNewFieldSplitter_EmptyField(t *testing.T) {
	_, err := NewFieldSplitter("", ",", "", 0)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewFieldSplitter_EmptyDelimiter(t *testing.T) {
	_, err := NewFieldSplitter("host", "", "", 0)
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNewFieldSplitter_DefaultPrefix(t *testing.T) {
	fs, err := NewFieldSplitter("tags", ",", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fs.prefix != "tags" {
		t.Errorf("expected prefix 'tags', got %q", fs.prefix)
	}
}

func TestFieldSplitter_SplitsField(t *testing.T) {
	fs, _ := NewFieldSplitter("tags", ",", "tag", 0)
	line := makeSplitLine(map[string]interface{}{"tags": "a,b,c"})
	out := fs.Process(line)

	for i, want := range []string{"a", "b", "c"} {
		key := "tag_" + string(rune('0'+i))
		got, ok := out.Fields[key]
		if !ok {
			t.Errorf("missing field %s", key)
			continue
		}
		if got != want {
			t.Errorf("field %s: got %q, want %q", key, got, want)
		}
	}
	// Original field preserved.
	if out.Fields["tags"] != "a,b,c" {
		t.Errorf("original field 'tags' should be preserved")
	}
}

func TestFieldSplitter_MaxParts(t *testing.T) {
	fs, _ := NewFieldSplitter("path", "/", "seg", 2)
	line := makeSplitLine(map[string]interface{}{"path": "a/b/c/d"})
	out := fs.Process(line)

	if _, ok := out.Fields["seg_2"]; ok {
		t.Error("expected at most 2 parts (seg_0, seg_1)")
	}
	if out.Fields["seg_0"] != "a" {
		t.Errorf("seg_0: got %v, want 'a'", out.Fields["seg_0"])
	}
	// With SplitN(n=2) the remainder goes into seg_1.
	if out.Fields["seg_1"] != "b/c/d" {
		t.Errorf("seg_1: got %v, want 'b/c/d'", out.Fields["seg_1"])
	}
}

func TestFieldSplitter_MissingField_Passthrough(t *testing.T) {
	fs, _ := NewFieldSplitter("tags", ",", "tag", 0)
	line := makeSplitLine(map[string]interface{}{"level": "info"})
	out := fs.Process(line)
	if len(out.Fields) != 1 {
		t.Errorf("expected no new fields, got %v", out.Fields)
	}
}

func TestFieldSplitter_NonStringField_Passthrough(t *testing.T) {
	fs, _ := NewFieldSplitter("count", ",", "part", 0)
	line := makeSplitLine(map[string]interface{}{"count": 42})
	out := fs.Process(line)
	if _, ok := out.Fields["part_0"]; ok {
		t.Error("should not split non-string field")
	}
}
