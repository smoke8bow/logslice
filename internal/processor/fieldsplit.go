package processor

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// FieldSplitter splits a string field value by a delimiter and stores the
// resulting parts as new fields with indexed names (e.g. field_0, field_1).
type FieldSplitter struct {
	field     string
	delimiter string
	prefix    string
	maxParts  int
}

// NewFieldSplitter creates a FieldSplitter that splits the given field by
// delimiter into at most maxParts sub-fields named "<prefix>_0", "<prefix>_1",
// etc. If prefix is empty, the source field name is used as the prefix.
// maxParts <= 0 means no limit.
func NewFieldSplitter(field, delimiter, prefix string, maxParts int) (*FieldSplitter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldSplitter: field name must not be empty")
	}
	if delimiter == "" {
		return nil, fmt.Errorf("fieldSplitter: delimiter must not be empty")
	}
	if prefix == "" {
		prefix = field
	}
	return &FieldSplitter{
		field:     field,
		delimiter: delimiter,
		prefix:    prefix,
		maxParts:  maxParts,
	}, nil
}

// Process splits the configured field in the log line and injects the resulting
// parts as new fields. The original field is preserved. Lines that do not
// contain the field, or whose field value is not a string, are passed through
// unchanged.
func (fs *FieldSplitter) Process(line parser.LogLine) parser.LogLine {
	val, ok := line.Fields[fs.field]
	if !ok {
		return line
	}
	str, ok := val.(string)
	if !ok {
		return line
	}

	n := -1
	if fs.maxParts > 0 {
		n = fs.maxParts
	}
	parts := strings.SplitN(str, fs.delimiter, n)

	// Copy fields map so we do not mutate the original.
	newFields := make(map[string]interface{}, len(line.Fields)+len(parts))
	for k, v := range line.Fields {
		newFields[k] = v
	}
	for i, p := range parts {
		newFields[fmt.Sprintf("%s_%d", fs.prefix, i)] = p
	}

	line.Fields = newFields
	return line
}
