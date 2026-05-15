package processor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldFormatter applies a sprintf-style format string to a named field value.
type FieldFormatter struct {
	field  string
	format string
}

// NewFieldFormatter creates a FieldFormatter that rewrites the value of field
// using the provided Go format string (e.g. "%.2f", "%08d", "%q").
func NewFieldFormatter(field, format string) (*FieldFormatter, error) {
	if strings.TrimSpace(field) == "" {
		return nil, fmt.Errorf("fieldformat: field name must not be empty")
	}
	if strings.TrimSpace(format) == "" {
		return nil, fmt.Errorf("fieldformat: format string must not be empty")
	}
	// Validate the format string is usable by doing a dry run with a zero value.
	result := fmt.Sprintf(format, 0)
	if strings.Contains(result, "%!(EXTRA") {
		return nil, fmt.Errorf("fieldformat: format string %q appears invalid", format)
	}
	return &FieldFormatter{field: field, format: format}, nil
}

// Process rewrites the target field of line using the configured format string.
// If the field is absent or the line is unstructured, the line is passed through unchanged.
func (f *FieldFormatter) Process(line parser.LogLine) parser.LogLine {
	if line.Fields == nil {
		return line
	}
	raw, ok := line.Fields[f.field]
	if !ok {
		return line
	}

	formatted := applyFormat(f.format, raw)

	newFields := make(map[string]string, len(line.Fields))
	for k, v := range line.Fields {
		newFields[k] = v
	}
	newFields[f.field] = formatted
	line.Fields = newFields
	return line
}

// applyFormat attempts numeric formatting first; falls back to string formatting.
func applyFormat(format, value string) string {
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return fmt.Sprintf(format, i)
	}
	if fl, err := strconv.ParseFloat(value, 64); err == nil {
		return fmt.Sprintf(format, fl)
	}
	return fmt.Sprintf(format, value)
}
