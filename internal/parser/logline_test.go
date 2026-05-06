package parser

import (
	"testing"
	"time"
)

func TestDetect_JSON(t *testing.T) {
	if got := detect(`{"msg":"hello"}`); got != FormatJSON {
		t.Fatalf("expected FormatJSON, got %v", got)
	}
}

func TestDetect_Logfmt(t *testing.T) {
	if got := detect(`level=info msg=hello`); got != FormatLogfmt {
		t.Fatalf("expected FormatLogfmt, got %v", got)
	}
}

func TestDetect_Plain(t *testing.T) {
	if got := detect(`plain log line`); got != FormatPlain {
		t.Fatalf("expected FormatPlain, got %v", got)
	}
}

func TestParser_ParseJSON(t *testing.T) {
	p := NewParser(FormatUnknown, nil)
	raw := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"started"}`
	line := p.Parse(raw)
	if line.Format != FormatJSON {
		t.Errorf("expected FormatJSON")
	}
	if line.Level != "info" {
		t.Errorf("expected level info, got %q", line.Level)
	}
	if line.Message != "started" {
		t.Errorf("expected message 'started', got %q", line.Message)
	}
	if line.Timestamp.IsZero() {
		t.Errorf("expected non-zero timestamp")
	}
}

func TestParser_ParseLogfmt(t *testing.T) {
	p := NewParser(FormatUnknown, nil)
	raw := `time=2024-01-15T10:00:00Z level=warn msg="disk full"`
	line := p.Parse(raw)
	if line.Format != FormatLogfmt {
		t.Errorf("expected FormatLogfmt")
	}
	if line.Level != "warn" {
		t.Errorf("expected level warn, got %q", line.Level)
	}
	if line.Message != "disk full" {
		t.Errorf("expected message 'disk full', got %q", line.Message)
	}
	if line.Timestamp.IsZero() {
		t.Errorf("expected non-zero timestamp")
	}
}

func TestParser_ParsePlain(t *testing.T) {
	p := NewParser(FormatPlain, nil)
	raw := "2024-01-15 plain text log"
	line := p.Parse(raw)
	if line.Format != FormatPlain {
		t.Errorf("expected FormatPlain")
	}
	if line.Message != raw {
		t.Errorf("expected raw message")
	}
	if !line.Timestamp.Equal(time.Time{}) {
		t.Errorf("expected zero timestamp for plain format")
	}
}

func TestParser_ExplicitFormat(t *testing.T) {
	p := NewParser(FormatJSON, nil)
	raw := `{"msg":"explicit"}`
	line := p.Parse(raw)
	if line.Format != FormatJSON {
		t.Errorf("expected FormatJSON")
	}
}
