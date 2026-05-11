package processor

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// FieldFilter filters log lines based on a field key-value match.
// It supports exact string matching and regexp matching on structured fields.
type FieldFilter struct {
	field   string
	exact   string
	pattern *regexp.Regexp
}

// NewFieldFilter creates a FieldFilter that matches lines where the given
// field equals value (exact) or matches valueRegexp (compiled regex).
// Exactly one of value or valueRegexp must be non-empty.
func NewFieldFilter(field, value, valueRegexp string) (*FieldFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("field name must not be empty")
	}
	if value == "" && valueRegexp == "" {
		return nil, fmt.Errorf("one of value or valueRegexp must be provided")
	}
	if value != "" && valueRegexp != "" {
		return nil, fmt.Errorf("only one of value or valueRegexp may be provided")
	}

	ff := &FieldFilter{field: field, exact: value}
	if valueRegexp != "" {
		re, err := regexp.Compile(valueRegexp)
		if err != nil {
			return nil, fmt.Errorf("invalid valueRegexp: %w", err)
		}
		ff.pattern = re
	}
	return ff, nil
}

// Keep returns true if the log line should be kept.
// Lines without structured fields are dropped when a field filter is active.
func (f *FieldFilter) Keep(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Fields[f.field]
	if !ok {
		return false
	}
	if f.pattern != nil {
		return f.pattern.MatchString(v)
	}
	return v == f.exact
}
