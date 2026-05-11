package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeFieldLogLine(fields map[string]string) *parser.LogLine {
	return &parser.LogLine{
		Raw:    "test line",
		Fields: fields,
	}
}

func TestNewFieldFilter_EmptyField(t *testing.T) {
	_, err := NewFieldFilter("", "val", "")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewFieldFilter_NeitherValueProvided(t *testing.T) {
	_, err := NewFieldFilter("level", "", "")
	if err == nil {
		t.Fatal("expected error when neither value nor valueRegexp provided")
	}
}

func TestNewFieldFilter_BothValuesProvided(t *testing.T) {
	_, err := NewFieldFilter("level", "info", "info")
	if err == nil {
		t.Fatal("expected error when both value and valueRegexp provided")
	}
}

func TestNewFieldFilter_InvalidRegexp(t *testing.T) {
	_, err := NewFieldFilter("level", "", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regexp")
	}
}

func TestFieldFilter_ExactMatch_Keeps(t *testing.T) {
	f, _ := NewFieldFilter("level", "error", "")
	line := makeFieldLogLine(map[string]string{"level": "error"})
	if !f.Keep(line) {
		t.Error("expected line to be kept")
	}
}

func TestFieldFilter_ExactMatch_Drops(t *testing.T) {
	f, _ := NewFieldFilter("level", "error", "")
	line := makeFieldLogLine(map[string]string{"level": "info"})
	if f.Keep(line) {
		t.Error("expected line to be dropped")
	}
}

func TestFieldFilter_MissingField_Drops(t *testing.T) {
	f, _ := NewFieldFilter("level", "error", "")
	line := makeFieldLogLine(map[string]string{"msg": "hello"})
	if f.Keep(line) {
		t.Error("expected line without field to be dropped")
	}
}

func TestFieldFilter_RegexpMatch_Keeps(t *testing.T) {
	f, _ := NewFieldFilter("service", "", "^auth.*")
	line := makeFieldLogLine(map[string]string{"service": "auth-service"})
	if !f.Keep(line) {
		t.Error("expected line to be kept by regexp")
	}
}

func TestFieldFilter_RegexpMatch_Drops(t *testing.T) {
	f, _ := NewFieldFilter("service", "", "^auth.*")
	line := makeFieldLogLine(map[string]string{"service": "payment-service"})
	if f.Keep(line) {
		t.Error("expected line to be dropped by regexp")
	}
}

func TestFieldFilter_NilLine_Drops(t *testing.T) {
	f, _ := NewFieldFilter("level", "error", "")
	if f.Keep(nil) {
		t.Error("expected nil line to be dropped")
	}
}
