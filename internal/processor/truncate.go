package processor

import (
	"fmt"
	"unicode/utf8"

	"github.com/user/logslice/internal/parser"
)

// Truncator trims log line messages or raw content that exceed a maximum byte length.
type Truncator struct {
	maxBytes int
	field    string
	suffix   string
}

// NewTruncator creates a Truncator that limits the given field (or raw line if field is empty)
// to maxBytes bytes. A suffix such as "..." is appended when truncation occurs.
func NewTruncator(maxBytes int, field, suffix string) (*Truncator, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("truncate: maxBytes must be positive, got %d", maxBytes)
	}
	if suffix == "" {
		suffix = "..."
	}
	if len(suffix) >= maxBytes {
		return nil, fmt.Errorf("truncate: suffix %q is >= maxBytes %d", suffix, maxBytes)
	}
	return &Truncator{maxBytes: maxBytes, field: field, suffix: suffix}, nil
}

// Apply truncates the relevant content of the log line and returns the (possibly modified) line.
func (t *Truncator) Apply(line parser.LogLine) (parser.LogLine, bool) {
	if t.field == "" {
		// Truncate raw line
		if len(line.Raw) <= t.maxBytes {
			return line, true
		}
		line.Raw = truncateString(line.Raw, t.maxBytes, t.suffix)
		return line, true
	}

	// Truncate a named field
	val, ok := line.Fields[t.field]
	if !ok {
		return line, true
	}
	s, ok2 := val.(string)
	if !ok2 {
		return line, true
	}
	if len(s) <= t.maxBytes {
		return line, true
	}
	if line.Fields == nil {
		line.Fields = make(map[string]interface{})
	}
	line.Fields[t.field] = truncateString(s, t.maxBytes, t.suffix)
	return line, true
}

// truncateString cuts s to at most maxBytes bytes (respecting UTF-8 boundaries) and appends suffix.
func truncateString(s string, maxBytes int, suffix string) string {
	cutAt := maxBytes - len(suffix)
	if cutAt <= 0 {
		return suffix
	}
	// Walk back to a valid UTF-8 boundary
	for cutAt > 0 && !utf8.RuneStart(s[cutAt]) {
		cutAt--
	}
	return s[:cutAt] + suffix
}
