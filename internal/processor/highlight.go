package processor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
)

// Highlighter applies ANSI color highlights to matching patterns in log lines.
type Highlighter struct {
	patterns []*regexp.Regexp
	color    string
}

// NewHighlighter creates a Highlighter that colorizes occurrences of any of
// the provided regex patterns. color must be one of: red, yellow, cyan, green.
func NewHighlighter(patterns []string, color string) (*Highlighter, error) {
	if len(patterns) == 0 {
		return nil, fmt.Errorf("highlight: at least one pattern is required")
	}

	ansi, err := resolveColor(color)
	if err != nil {
		return nil, err
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("highlight: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}

	return &Highlighter{patterns: compiled, color: ansi}, nil
}

// Apply colorizes matched substrings within line.Raw and returns the modified line.
func (h *Highlighter) Apply(line parser.LogLine) parser.LogLine {
	raw := line.Raw
	for _, re := range h.patterns {
		raw = re.ReplaceAllStringFunc(raw, func(match string) string {
			return h.color + match + colorReset
		})
	}
	line.Raw = raw
	return line
}

// resolveColor maps a color name to its ANSI escape sequence.
func resolveColor(name string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "red":
		return colorRed, nil
	case "yellow":
		return colorYellow, nil
	case "cyan":
		return colorCyan, nil
	case "green":
		return colorGreen, nil
	case "":
		return colorCyan, nil // default
	default:
		return "", fmt.Errorf("highlight: unknown color %q (want: red, yellow, cyan, green)", name)
	}
}
