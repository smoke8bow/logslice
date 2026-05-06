package output

import (
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"", FormatRaw},
		{"raw", FormatRaw},
		{"RAW", FormatRaw},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
		{"tsv", FormatTSV},
		{"TSV", FormatTSV},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, got)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestFormatLine_Raw(t *testing.T) {
	result := FormatLine("my log line", 1, FormatRaw)
	if result != "my log line" {
		t.Errorf("expected raw line, got %q", result)
	}
}

func TestFormatLine_JSON(t *testing.T) {
	result := FormatLine("hello", 3, FormatJSON)
	if !strings.Contains(result, `"n":3`) {
		t.Errorf("expected line number in JSON output, got %q", result)
	}
	if !strings.Contains(result, `"line":"hello"`) {
		t.Errorf("expected line content in JSON output, got %q", result)
	}
}

func TestFormatLine_JSON_EscapesQuotes(t *testing.T) {
	result := FormatLine(`say "hi"`, 1, FormatJSON)
	if strings.Contains(result, `say "hi"`) {
		t.Errorf("expected escaped quotes in JSON output, got %q", result)
	}
}

func TestFormatLine_TSV(t *testing.T) {
	result := FormatLine("tab line", 7, FormatTSV)
	expected := "7\ttab line"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
