package processor

import (
	"errors"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldRenamer renames fields in structured log lines.
type FieldRenamer struct {
	mappings map[string]string
}

// NewFieldRenamer creates a FieldRenamer from a slice of "old=new" mapping strings.
// Returns an error if any mapping is malformed or contains blank field names.
func NewFieldRenamer(mappings []string) (*FieldRenamer, error) {
	if len(mappings) == 0 {
		return nil, errors.New("rename: at least one mapping is required")
	}
	m := make(map[string]string, len(mappings))
	for _, entry := range mappings {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("rename: invalid mapping format, expected old=new: " + entry)
		}
		old := strings.TrimSpace(parts[0])
		new := strings.TrimSpace(parts[1])
		if old == "" {
			return nil, errors.New("rename: source field name must not be blank")
		}
		if new == "" {
			return nil, errors.New("rename: destination field name must not be blank")
		}
		m[old] = new
	}
	return &FieldRenamer{mappings: m}, nil
}

// Process renames configured fields in the log line's Fields map.
// Unstructured (plain) lines are passed through unchanged.
func (r *FieldRenamer) Process(line *parser.LogLine) *parser.LogLine {
	if line == nil || len(line.Fields) == 0 {
		return line
	}
	for old, newName := range r.mappings {
		if val, ok := line.Fields[old]; ok {
			delete(line.Fields, old)
			line.Fields[newName] = val
		}
	}
	return line
}
