package processor

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelFatal,
	"crit":  LevelFatal,
}

// ParseLevel converts a string to a Level. Returns LevelUnknown if unrecognised.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// LevelFilter discards log lines whose severity is below the minimum level.
type LevelFilter struct {
	min Level
}

// NewLevelFilter creates a LevelFilter that keeps lines at or above minLevel.
// minLevel must be one of: debug, info, warn, error, fatal.
func NewLevelFilter(minLevel string) (*LevelFilter, error) {
	l := ParseLevel(minLevel)
	if l == LevelUnknown {
		return nil, &ErrUnknownLevel{Raw: minLevel}
	}
	return &LevelFilter{min: l}, nil
}

// Keep returns true when the line's level is at or above the minimum.
func (f *LevelFilter) Keep(line *parser.LogLine) bool {
	raw := line.Fields["level"]
	if raw == "" {
		raw = line.Fields["lvl"]
	}
	if raw == "" {
		// No level field — keep the line to avoid false negatives.
		return true
	}
	l := ParseLevel(raw)
	if l == LevelUnknown {
		return true
	}
	return l >= f.min
}

// ErrUnknownLevel is returned when an unrecognised level string is supplied.
type ErrUnknownLevel struct {
	Raw string
}

func (e *ErrUnknownLevel) Error() string {
	return "unknown log level: " + e.Raw
}
