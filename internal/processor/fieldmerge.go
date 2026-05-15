package processor

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldMerger combines multiple source fields into a single destination field
// using a configurable separator.
type FieldMerger struct {
	sources []string
	dest    string
	sep     string
	remove  bool
}

// NewFieldMerger creates a FieldMerger that concatenates the given source fields
// into dest using sep as the separator. If remove is true, source fields are
// deleted after merging. At least two source fields must be provided.
func NewFieldMerger(sources []string, dest, sep string, remove bool) (*FieldMerger, error) {
	if len(sources) < 2 {
		return nil, fmt.Errorf("fieldmerge: at least two source fields required")
	}
	for i, s := range sources {
		if strings.TrimSpace(s) == "" {
			return nil, fmt.Errorf("fieldmerge: source field at index %d is blank", i)
		}
	}
	if strings.TrimSpace(dest) == "" {
		return nil, fmt.Errorf("fieldmerge: destination field must not be blank")
	}
	return &FieldMerger{
		sources: sources,
		dest:    dest,
		sep:     sep,
		remove:  remove,
	}, nil
}

// Process merges the source fields of line into the destination field.
// Lines without structured fields are returned unchanged.
func (fm *FieldMerger) Process(line parser.LogLine) (parser.LogLine, bool) {
	if len(line.Fields) == 0 {
		return line, true
	}

	parts := make([]string, 0, len(fm.sources))
	for _, src := range fm.sources {
		if v, ok := line.Fields[src]; ok {
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}

	if len(parts) == 0 {
		return line, true
	}

	merged := strings.Join(parts, fm.sep)

	// Clone fields map to avoid mutating the original.
	newFields := make(map[string]interface{}, len(line.Fields)+1)
	for k, v := range line.Fields {
		newFields[k] = v
	}
	newFields[fm.dest] = merged

	if fm.remove {
		for _, src := range fm.sources {
			if src != fm.dest {
				delete(newFields, src)
			}
		}
	}

	line.Fields = newFields
	return line, true
}
