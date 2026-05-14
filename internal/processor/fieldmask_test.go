package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeMaskLine(fields map[string]string) parser.LogLine {
	return parser.LogLine{
		Raw:    "level=info user=alice token=secret",
		Fields: fields,
	}
}

func TestNewFieldMask_EmptyFields(t *testing.T) {
	_, err := NewFieldMask([]string{}, "")
	if err == nil {
		t.Fatal("expected error for empty fields list")
	}
}

func TestNewFieldMask_BlankFieldName(t *testing.T) {
	_, err := NewFieldMask([]string{"  "}, "REDACTED")
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestFieldMask_RemovesField(t *testing.T) {
	fm, err := NewFieldMask([]string{"token"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeMaskLine(map[string]string{"level": "info", "user": "alice", "token": "secret"})
	out := fm.Apply(line)
	if _, ok := out.Fields["token"]; ok {
		t.Error("expected 'token' field to be removed")
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("expected 'user' to be preserved, got %q", out.Fields["user"])
	}
}

func TestFieldMask_RedactsField(t *testing.T) {
	fm, err := NewFieldMask([]string{"token", "password"}, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeMaskLine(map[string]string{"level": "info", "token": "abc123", "password": "hunter2"})
	out := fm.Apply(line)
	if out.Fields["token"] != "***" {
		t.Errorf("expected token redacted to ***, got %q", out.Fields["token"])
	}
	if out.Fields["password"] != "***" {
		t.Errorf("expected password redacted to ***, got %q", out.Fields["password"])
	}
	if out.Fields["level"] != "info" {
		t.Errorf("expected level preserved, got %q", out.Fields["level"])
	}
}

func TestFieldMask_NoFieldsInLine(t *testing.T) {
	fm, err := NewFieldMask([]string{"token"}, "REDACTED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := parser.LogLine{Raw: "plain log line", Fields: nil}
	out := fm.Apply(line)
	if out.Raw != line.Raw {
		t.Errorf("expected raw to be unchanged, got %q", out.Raw)
	}
}

func TestFieldMask_MultipleFieldsRemoved(t *testing.T) {
	fm, err := NewFieldMask([]string{"token", "user"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeMaskLine(map[string]string{"level": "warn", "user": "bob", "token": "xyz"})
	out := fm.Apply(line)
	if _, ok := out.Fields["token"]; ok {
		t.Error("expected 'token' to be removed")
	}
	if _, ok := out.Fields["user"]; ok {
		t.Error("expected 'user' to be removed")
	}
	if out.Fields["level"] != "warn" {
		t.Errorf("expected 'level' preserved, got %q", out.Fields["level"])
	}
}
