package processor

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeTransformLine(fields map[string]interface{}) parser.LogLine {
	return parser.LogLine{
		Raw:    "raw log line",
		Fields: fields,
	}
}

func TestNewTransformer_EmptyField(t *testing.T) {
	_, err := NewTransformer("", "upper")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTransformer_UnknownTransform(t *testing.T) {
	_, err := NewTransformer("msg", "rot13")
	if err == nil {
		t.Fatal("expected error for unknown transform")
	}
}

func TestNewTransformer_ValidTransforms(t *testing.T) {
	for _, name := range []string{"upper", "lower", "trim", "urlencode"} {
		_, err := NewTransformer("field", name)
		if err != nil {
			t.Errorf("unexpected error for transform %q: %v", name, err)
		}
	}
}

func TestTransformer_UpperCase(t *testing.T) {
	tr, _ := NewTransformer("msg", "upper")
	line := makeTransformLine(map[string]interface{}{"msg": "hello world"})
	out := tr.Apply(line)
	if out.Fields["msg"] != "HELLO WORLD" {
		t.Errorf("expected HELLO WORLD, got %v", out.Fields["msg"])
	}
}

func TestTransformer_LowerCase(t *testing.T) {
	tr, _ := NewTransformer("level", "lower")
	line := makeTransformLine(map[string]interface{}{"level": "ERROR"})
	out := tr.Apply(line)
	if out.Fields["level"] != "error" {
		t.Errorf("expected error, got %v", out.Fields["level"])
	}
}

func TestTransformer_Trim(t *testing.T) {
	tr, _ := NewTransformer("msg", "trim")
	line := makeTransformLine(map[string]interface{}{"msg": "  padded  "})
	out := tr.Apply(line)
	if out.Fields["msg"] != "padded" {
		t.Errorf("expected padded, got %q", out.Fields["msg"])
	}
}

func TestTransformer_URLEncode(t *testing.T) {
	tr, _ := NewTransformer("path", "urlencode")
	line := makeTransformLine(map[string]interface{}{"path": "/foo bar"})
	out := tr.Apply(line)
	if out.Fields["path"] != "%2Ffoo%20bar" {
		t.Errorf("unexpected urlencode result: %v", out.Fields["path"])
	}
}

func TestTransformer_MissingField_NoOp(t *testing.T) {
	tr, _ := NewTransformer("missing", "upper")
	line := makeTransformLine(map[string]interface{}{"msg": "hello"})
	out := tr.Apply(line)
	if _, ok := out.Fields["missing"]; ok {
		t.Error("expected missing field to remain absent")
	}
	if out.Fields["msg"] != "hello" {
		t.Error("expected unrelated field to be unchanged")
	}
}

func TestTransformer_NonStringField_NoOp(t *testing.T) {
	tr, _ := NewTransformer("count", "upper")
	line := makeTransformLine(map[string]interface{}{"count": 42})
	out := tr.Apply(line)
	if out.Fields["count"] != 42 {
		t.Errorf("expected non-string field to be unchanged, got %v", out.Fields["count"])
	}
}

func TestTransformer_DoesNotMutateOriginal(t *testing.T) {
	tr, _ := NewTransformer("msg", "upper")
	origFields := map[string]interface{}{"msg": "hello"}
	line := makeTransformLine(origFields)
	tr.Apply(line)
	if origFields["msg"] != "hello" {
		t.Error("original fields map was mutated")
	}
}

// TestTransformer_RawPreserved verifies that Apply does not modify the Raw
// field of the log line regardless of which transform is applied.
func TestTransformer_RawPreserved(t *testing.T) {
	for _, name := range []string{"upper", "lower", "trim", "urlencode"} {
		tr, _ := NewTransformer("msg", name)
		line := makeTransformLine(map[string]interface{}{"msg": "hello"})
		out := tr.Apply(line)
		if out.Raw != line.Raw {
			t.Errorf("transform %q modified Raw field: got %q, want %q", name, out.Raw, line.Raw)
		}
	}
}
