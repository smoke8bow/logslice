package parser

import (
	"strings"
	"time"
)

// Format represents the detected or configured log format.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatLogfmt
	FormatPlain
)

// LogLine represents a parsed log line with optional structured fields.
type LogLine struct {
	Raw       string
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
	Format    Format
}

// Parser parses raw log lines into LogLine structs.
type Parser struct {
	detectFormat bool
	format       Format
	timeFields   []string
}

// NewParser creates a Parser. If format is FormatUnknown, auto-detection is used.
func NewParser(format Format, timeFields []string) *Parser {
	if len(timeFields) == 0 {
		timeFields = []string{"time", "ts", "timestamp", "@timestamp"}
	}
	return &Parser{
		detectFormat: format == FormatUnknown,
		format:       format,
		timeFields:   timeFields,
	}
}

// Parse parses a raw string into a LogLine.
func (p *Parser) Parse(raw string) LogLine {
	line := LogLine{Raw: raw, Fields: make(map[string]string)}
	fmt := p.format
	if p.detectFormat {
		fmt = detect(raw)
	}
	line.Format = fmt
	switch fmt {
	case FormatJSON:
		parseJSON(raw, &line, p.timeFields)
	case FormatLogfmt:
		parseLogfmt(raw, &line, p.timeFields)
	default:
		line.Message = raw
	}
	return line
}

// detect guesses the format of a raw log line.
func detect(raw string) Format {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return FormatJSON
	}
	if strings.Contains(trimmed, "=") {
		return FormatLogfmt
	}
	return FormatPlain
}
