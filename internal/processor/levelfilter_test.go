package processor

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeLogLine(level string) *parser.LogLine {
	fields := map[string]string{}
	if level != "" {
		fields["level"] = level
	}
	return &parser.LogLine{Fields: fields, Raw: "test line"}
}

func TestParseLevel_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"err", LevelError},
		{"fatal", LevelFatal},
		{"crit", LevelFatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			if got := ParseLevel(tc.input); got != tc.want {
				t.Errorf("ParseLevel(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	if got := ParseLevel("trace"); got != LevelUnknown {
		t.Errorf("expected LevelUnknown, got %d", got)
	}
}

func TestNewLevelFilter_InvalidLevel(t *testing.T) {
	_, err := NewLevelFilter("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestLevelFilter_KeepsAtOrAboveMin(t *testing.T) {
	f, err := NewLevelFilter("warn")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		level string
		want  bool
	}{
		{"debug", false},
		{"info", false},
		{"warn", true},
		{"error", true},
		{"fatal", true},
	}
	for _, tc := range cases {
		t.Run(tc.level, func(t *testing.T) {
			line := makeLogLine(tc.level)
			if got := f.Keep(line); got != tc.want {
				t.Errorf("Keep(%q) = %v, want %v", tc.level, got, tc.want)
			}
		})
	}
}

func TestLevelFilter_NoLevelField_Kept(t *testing.T) {
	f, _ := NewLevelFilter("error")
	line := makeLogLine("")
	if !f.Keep(line) {
		t.Error("line with no level field should be kept")
	}
}

func TestLevelFilter_UnknownLevelField_Kept(t *testing.T) {
	f, _ := NewLevelFilter("error")
	line := makeLogLine("trace")
	if !f.Keep(line) {
		t.Error("line with unrecognised level should be kept")
	}
}
