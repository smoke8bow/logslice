package parser

import (
	"strings"
	"time"
)

// commonTimeLayouts lists timestamp formats tried when parsing plain log lines.
var commonTimeLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"Jan 02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// parsePlain attempts to extract a timestamp and message from an unstructured
// log line. It scans the first few whitespace-delimited tokens looking for a
// recognisable timestamp. The remainder of the line is stored as the message.
func parsePlain(line string) LogLine {
	ll := LogLine{
		Raw:    line,
		Fields: make(map[string]string),
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return ll
	}

	// Try progressively longer token prefixes as timestamps (up to 4 tokens).
	tokens := strings.Fields(trimmed)
	for n := 1; n <= 4 && n <= len(tokens); n++ {
		candidate := strings.Join(tokens[:n], " ")
		if t, ok := tryParseTime(candidate); ok {
			ll.Timestamp = t
			if n < len(tokens) {
				ll.Message = strings.TrimSpace(strings.Join(tokens[n:], " "))
			}
			ll.Fields["msg"] = ll.Message
			return ll
		}
	}

	// No timestamp found — treat the whole line as the message.
	ll.Message = trimmed
	ll.Fields["msg"] = trimmed
	return ll
}

// tryParseTime attempts to parse s against all known layouts.
func tryParseTime(s string) (time.Time, bool) {
	for _, layout := range commonTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
