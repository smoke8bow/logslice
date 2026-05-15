package processor

import (
	"errors"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldAdder injects static key=value pairs into structured log lines.
type FieldAdder struct {
	fields map[string]string
	overwrite bool
}

// NewFieldAdder creates a FieldAdder from a slice of "key=value" strings.
// If overwrite is true, existing fields with the same name will be replaced.
func NewFieldAdder(pairs []string, overwrite bool) (*FieldAdder, error) {
	if len(pairs) == 0 {
		return nil, errors.New("addfield: at least one key=value pair is required")
	}
	fields := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("addfield: invalid pair format, expected key=value: " + p)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, errors.New("addfield: field name must not be blank")
		}
		fields[key] = val
	}
	return &FieldAdder{fields: fields, overwrite: overwrite}, nil
}

// Process injects the configured fields into the log line.
// Lines without an existing Fields map will have one initialised.
// Unstructured lines gain a Fields map but their Raw value is preserved.
func (a *FieldAdder) Process(line *parser.LogLine) *parser.LogLine {
	if line == nil {
		return nil
	}
	if line.Fields == nil {
		line.Fields = make(map[string]string, len(a.fields))
	}
	for k, v := range a.fields {
		if _, exists := line.Fields[k]; exists && !a.overwrite {
			continue
		}
		line.Fields[k] = v
	}
	return line
}
