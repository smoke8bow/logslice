package processor

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/user/logslice/internal/parser"
)

// Templater rewrites the Raw field of a LogLine using a Go text/template.
// Fields from the parsed log line are available as .Fields["key"],
// and .Raw contains the original raw line.
type Templater struct {
	tmpl *template.Template
}

type templateData struct {
	Raw    string
	Fields map[string]string
}

// NewTemplater compiles tmplStr and returns a Templater.
// Returns an error if the template is empty or fails to parse.
func NewTemplater(tmplStr string) (*Templater, error) {
	if strings.TrimSpace(tmplStr) == "" {
		return nil, fmt.Errorf("template: template string must not be empty")
	}
	tmpl, err := template.New("logline").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template: failed to compile: %w", err)
	}
	return &Templater{tmpl: tmpl}, nil
}

// Apply executes the template against the LogLine and updates Raw.
// If execution fails the original line is returned unchanged.
func (t *Templater) Apply(line parser.LogLine) parser.LogLine {
	data := templateData{
		Raw:    line.Raw,
		Fields: line.Fields,
	}
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return line
	}
	line.Raw = buf.String()
	return line
}
