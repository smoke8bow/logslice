package parser

import "testing"

func TestSplitLogfmt_Simple(t *testing.T) {
	m := splitLogfmt(`key=value foo=bar`)
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %q", m["key"])
	}
	if m["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %q", m["foo"])
	}
}

func TestSplitLogfmt_QuotedValue(t *testing.T) {
	m := splitLogfmt(`msg="hello world" level=info`)
	if m["msg"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", m["msg"])
	}
	if m["level"] != "info" {
		t.Errorf("expected 'info', got %q", m["level"])
	}
}

func TestSplitLogfmt_EmptyInput(t *testing.T) {
	m := splitLogfmt("")
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestSplitLogfmt_NoEquals(t *testing.T) {
	m := splitLogfmt("no equals here")
	if len(m) != 0 {
		t.Errorf("expected empty map for no-equals input")
	}
}

func TestParseLogfmt_FieldExtraction(t *testing.T) {
	p := NewParser(FormatLogfmt, nil)
	line := p.Parse(`level=error msg="connection refused" host=db01`)
	if line.Level != "error" {
		t.Errorf("expected level error, got %q", line.Level)
	}
	if line.Message != "connection refused" {
		t.Errorf("expected message 'connection refused', got %q", line.Message)
	}
	if line.Fields["host"] != "db01" {
		t.Errorf("expected host=db01, got %q", line.Fields["host"])
	}
}
