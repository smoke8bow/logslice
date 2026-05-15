package processor

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeRenameLine(fields map[string]string) *parser.LogLine {
	return &parser.LogLine{
		Raw:    "ts=2024-01-01T00:00:00Z level=info msg=hello",
		Fields: fields,
		Time:   time.Time{},
	}
}

func TestNewFieldRenamer_EmptyMappings(t *testing.T) {
	_, err := NewFieldRenamer([]string{})
	if err == nil {
		t.Fatal("expected error for empty mappings")
	}
}

func TestNewFieldRenamer_InvalidFormat(t *testing.T) {
	_, err := NewFieldRenamer([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '=' in mapping")
	}
}

func TestNewFieldRenamer_BlankSource(t *testing.T) {
	_, err := NewFieldRenamer([]string{"=newname"})
	if err == nil {
		t.Fatal("expected error for blank source field")
	}
}

func TestNewFieldRenamer_BlankDestination(t *testing.T) {
	_, err := NewFieldRenamer([]string{"oldname="})
	if err == nil {
		t.Fatal("expected error for blank destination field")
	}
}

func TestFieldRenamer_RenamesSingleField(t *testing.T) {
	r, err := NewFieldRenamer([]string{"msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeRenameLine(map[string]string{"msg": "hello", "level": "info"})
	out := r.Process(line)
	if _, ok := out.Fields["msg"]; ok {
		t.Error("old field 'msg' should have been removed")
	}
	if out.Fields["message"] != "hello" {
		t.Errorf("expected Fields['message']='hello', got %q", out.Fields["message"])
	}
}

func TestFieldRenamer_RenamesMultipleFields(t *testing.T) {
	r, err := NewFieldRenamer([]string{"msg=message", "level=severity"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeRenameLine(map[string]string{"msg": "hello", "level": "info"})
	out := r.Process(line)
	if out.Fields["message"] != "hello" {
		t.Errorf("expected message='hello', got %q", out.Fields["message"])
	}
	if out.Fields["severity"] != "info" {
		t.Errorf("expected severity='info', got %q", out.Fields["severity"])
	}
}

func TestFieldRenamer_MissingFieldIsNoop(t *testing.T) {
	r, err := NewFieldRenamer([]string{"nonexistent=other"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeRenameLine(map[string]string{"msg": "hello"})
	out := r.Process(line)
	if out.Fields["msg"] != "hello" {
		t.Error("existing field should be untouched")
	}
}

func TestFieldRenamer_NilLine(t *testing.T) {
	r, _ := NewFieldRenamer([]string{"msg=message"})
	out := r.Process(nil)
	if out != nil {
		t.Error("expected nil output for nil input")
	}
}

func TestFieldRenamer_NoFields(t *testing.T) {
	r, _ := NewFieldRenamer([]string{"msg=message"})
	line := &parser.LogLine{Raw: "plain log line", Fields: nil}
	out := r.Process(line)
	if out.Raw != "plain log line" {
		t.Error("plain line should be passed through unchanged")
	}
}
