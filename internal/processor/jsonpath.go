package processor

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// JSONPathExtractor extracts a nested field from structured log lines
// using a dot-separated path (e.g. "request.headers.host").
type JSONPathExtractor struct {
	path   []string
	output string
}

// NewJSONPathExtractor creates a new extractor that reads from the dot-separated
// path and writes the result into the field named by outputField.
func NewJSONPathExtractor(dotPath, outputField string) (*JSONPathExtractor, error) {
	if dotPath == "" {
		return nil, fmt.Errorf("jsonpath: dot path must not be empty")
	}
	if outputField == "" {
		return nil, fmt.Errorf("jsonpath: output field must not be empty")
	}
	parts := strings.Split(dotPath, ".")
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("jsonpath: path segment must not be empty in %q", dotPath)
		}
	}
	return &JSONPathExtractor{path: parts, output: outputField}, nil
}

// Process extracts the nested value and stores it as a top-level field.
func (e *JSONPathExtractor) Process(line *parser.LogLine) (*parser.LogLine, error) {
	if line == nil {
		return nil, nil
	}
	if len(line.Fields) == 0 {
		return line, nil
	}
	val, ok := walkPath(line.Fields, e.path)
	if !ok {
		return line, nil
	}
	if line.Fields == nil {
		line.Fields = make(map[string]string)
	}
	line.Fields[e.output] = val
	return line, nil
}

// walkPath traverses a nested map structure represented as map[string]string
// using dot-separated keys encoded as "parent.child" within a single key.
func walkPath(fields map[string]string, path []string) (string, bool) {
	// Try exact dotted key first (some parsers store as-is).
	full := strings.Join(path, ".")
	if v, ok := fields[full]; ok {
		return v, true
	}
	// Try progressive prefix matching for parsers that flatten with dots.
	for i := len(path); i >= 1; i-- {
		prefix := strings.Join(path[:i], ".")
		if v, ok := fields[prefix]; ok && i == len(path) {
			return v, true
		}
		_ = prefix
	}
	return "", false
}
