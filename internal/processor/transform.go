package processor

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// TransformFunc is a function that transforms a field value.
type TransformFunc func(string) string

// Transformer applies value transformations to named fields in a log line.
type Transformer struct {
	field     string
	transform TransformFunc
	transformName string
}

// NewTransformer creates a Transformer that applies a named transform to the given field.
// Supported transform names: "upper", "lower", "trim", "urlencode".
func NewTransformer(field, transformName string) (*Transformer, error) {
	if field == "" {
		return nil, fmt.Errorf("transform: field name must not be empty")
	}
	fn, err := resolveTransform(transformName)
	if err != nil {
		return nil, err
	}
	return &Transformer{
		field:         field,
		transform:     fn,
		transformName: transformName,
	}, nil
}

// Apply transforms the specified field in the log line, returning a modified copy.
func (t *Transformer) Apply(line parser.LogLine) parser.LogLine {
	val, ok := line.Fields[t.field]
	if !ok {
		return line
	}
	str, ok := val.(string)
	if !ok {
		return line
	}
	result := t.transform(str)
	if line.Fields == nil {
		return line
	}
	newFields := make(map[string]interface{}, len(line.Fields))
	for k, v := range line.Fields {
		newFields[k] = v
	}
	newFields[t.field] = result
	line.Fields = newFields
	return line
}

func resolveTransform(name string) (TransformFunc, error) {
	switch strings.ToLower(name) {
	case "upper":
		return strings.ToUpper, nil
	case "lower":
		return strings.ToLower, nil
	case "trim":
		return strings.TrimSpace, nil
	case "urlencode":
		return urlEncodeValue, nil
	default:
		return nil, fmt.Errorf("transform: unknown transform %q (want: upper, lower, trim, urlencode)", name)
	}
}

func urlEncodeValue(s string) string {
	var b strings.Builder
	for _, c := range s {
		switch {
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '-' || c == '_' || c == '.' || c == '~':
			b.WriteRune(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}
