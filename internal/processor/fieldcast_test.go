package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeCastLine(fields map[string]interface{}) parser.LogLine {
	return parser.LogLine{
		Raw:    "raw log line",
		Fields: fields,
	}
}

func TestNewFieldCaster_EmptyField(t *testing.T) {
	_, err := NewFieldCaster("", CastInt)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewFieldCaster_UnknownType(t *testing.T) {
	_, err := NewFieldCaster("level", CastType("bytes"))
	if err == nil {
		t.Fatal("expected error for unknown cast type")
	}
}

func TestNewFieldCaster_Valid(t *testing.T) {
	for _, ct := range []CastType{CastInt, CastFloat, CastString, CastBool} {
		_, err := NewFieldCaster("field", ct)
		if err != nil {
			t.Fatalf("unexpected error for type %q: %v", ct, err)
		}
	}
}

func TestFieldCaster_CastsToInt(t *testing.T) {
	fc, _ := NewFieldCaster("count", CastInt)
	line := makeCastLine(map[string]interface{}{"count": "42"})
	out := fc.Process(line)
	val, ok := out.Fields["count"]
	if !ok {
		t.Fatal("field missing after cast")
	}
	if _, ok := val.(int64); !ok {
		t.Fatalf("expected int64, got %T", val)
	}
}

func TestFieldCaster_CastsToFloat(t *testing.T) {
	fc, _ := NewFieldCaster("ratio", CastFloat)
	line := makeCastLine(map[string]interface{}{"ratio": "3.14"})
	out := fc.Process(line)
	val := out.Fields["ratio"]
	if _, ok := val.(float64); !ok {
		t.Fatalf("expected float64, got %T", val)
	}
}

func TestFieldCaster_CastsToBool(t *testing.T) {
	fc, _ := NewFieldCaster("enabled", CastBool)
	line := makeCastLine(map[string]interface{}{"enabled": "true"})
	out := fc.Process(line)
	val := out.Fields["enabled"]
	if b, ok := val.(bool); !ok || !b {
		t.Fatalf("expected bool true, got %v (%T)", val, val)
	}
}

func TestFieldCaster_CastsToString(t *testing.T) {
	fc, _ := NewFieldCaster("code", CastString)
	line := makeCastLine(map[string]interface{}{"code": 404})
	out := fc.Process(line)
	val := out.Fields["code"]
	if _, ok := val.(string); !ok {
		t.Fatalf("expected string, got %T", val)
	}
}

func TestFieldCaster_MissingFieldPassthrough(t *testing.T) {
	fc, _ := NewFieldCaster("missing", CastInt)
	line := makeCastLine(map[string]interface{}{"other": "hello"})
	out := fc.Process(line)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestFieldCaster_InvalidValuePassthrough(t *testing.T) {
	fc, _ := NewFieldCaster("count", CastInt)
	line := makeCastLine(map[string]interface{}{"count": "not-a-number"})
	out := fc.Process(line)
	if out.Fields["count"] != "not-a-number" {
		t.Fatal("expected original value to be preserved on cast failure")
	}
}
