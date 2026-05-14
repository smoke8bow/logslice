package processor

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldMask redacts or removes specified fields from parsed log lines.
type FieldMask struct {
	fields  map[string]struct{}
	redact  string
	remove  bool
}

// NewFieldMask creates a FieldMask that either removes or redacts the given fields.
// If redactWith is empty, matched fields are removed entirely.
// Otherwise, their values are replaced with the redactWith string.
func NewFieldMask(fields []string, redactWith string) (*FieldMask, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("fieldmask: at least one field name is required")
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return nil, fmt.Errorf("fieldmask: field name must not be blank")
		}
		set[f] = struct{}{}
	}
	return &FieldMask{
		fields: set,
		redact: redactWith,
		remove: redactWith == "",
	}, nil
}

// Apply returns a new LogLine with the configured fields removed or redacted.
// The Raw field is updated to reflect the change for plain lines.
func (fm *FieldMask) Apply(line parser.LogLine) parser.LogLine {
	if len(line.Fields) == 0 {
		return line
	}
	newFields := make(map[string]string, len(line.Fields))
	for k, v := range line.Fields {
		if _, masked := fm.fields[k]; masked {
			if !fm.remove {
				newFields[k] = fm.redact
			}
			// if remove==true, skip the field entirely
		} else {
			newFields[k] = v
		}
	}
	line.Fields = newFields
	return line
}
