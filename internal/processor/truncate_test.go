package processor

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeTruncLine(raw string, fields map[string]interface{}) parser.LogLine {
	return parser.LogLine{Raw: raw, Fields: fields}
}

func TestNewTruncator_InvalidMaxBytes(t *testing.T) {
	_, err := NewTruncator(0, "", "...")
	if err == nil {
		t.Fatal("expected error for maxBytes=0")
	}
}

func TestNewTruncator_SuffixTooLong(t *testing.T) {
	_, err := NewTruncator(3, "", "....")
	if err == nil {
		t.Fatal("expected error when suffix >= maxBytes")
	}
}

func TestNewTruncator_DefaultSuffix(t *testing.T) {
	tr, err := NewTruncator(20, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.suffix != "..." {
		t.Errorf("expected default suffix '...', got %q", tr.suffix)
	}
}

func TestTruncator_RawLineShort(t *testing.T) {
	tr, _ := NewTruncator(50, "", "...")
	line := makeTruncLine("short line", nil)
	out, keep := tr.Apply(line)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if out.Raw != "short line" {
		t.Errorf("unexpected modification: %q", out.Raw)
	}
}

func TestTruncator_RawLineTruncated(t *testing.T) {
	tr, _ := NewTruncator(10, "", "...")
	line := makeTruncLine("this is a long log line", nil)
	out, keep := tr.Apply(line)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if len(out.Raw) > 10 {
		t.Errorf("expected truncated to <=10 bytes, got %d: %q", len(out.Raw), out.Raw)
	}
	if !strings.HasSuffix(out.Raw, "...") {
		t.Errorf("expected suffix '...', got %q", out.Raw)
	}
}

func TestTruncator_FieldTruncated(t *testing.T) {
	tr, _ := NewTruncator(8, "msg", "...")
	fields := map[string]interface{}{"msg": "hello world this is long"}
	line := makeTruncLine("", fields)
	out, keep := tr.Apply(line)
	if !keep {
		t.Fatal("expected keep=true")
	}
	val, _ := out.Fields["msg"].(string)
	if len(val) > 8 {
		t.Errorf("field not truncated: %q", val)
	}
	if !strings.HasSuffix(val, "...") {
		t.Errorf("expected suffix '...', got %q", val)
	}
}

func TestTruncator_FieldMissing(t *testing.T) {
	tr, _ := NewTruncator(8, "msg", "...")
	line := makeTruncLine("raw", map[string]interface{}{"level": "info"})
	out, keep := tr.Apply(line)
	if !keep {
		t.Fatal("expected keep=true")
	}
	if out.Raw != "raw" {
		t.Errorf("raw unexpectedly modified: %q", out.Raw)
	}
}

func TestTruncator_FieldNotString(t *testing.T) {
	tr, _ := NewTruncator(8, "count", "...")
	fields := map[string]interface{}{"count": 12345}
	line := makeTruncLine("", fields)
	out, _ := tr.Apply(line)
	if out.Fields["count"] != 12345 {
		t.Errorf("non-string field was modified")
	}
}
