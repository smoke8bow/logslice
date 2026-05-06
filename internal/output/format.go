package output

import (
	"fmt"
	"strings"
)

// Format represents an output format for log lines.
type Format int

const (
	// FormatRaw outputs lines as-is.
	FormatRaw Format = iota
	// FormatJSON wraps each line in a simple JSON envelope.
	FormatJSON
	// FormatTSV outputs lines with a tab-separated prefix (line number).
	FormatTSV
)

// ParseFormat converts a string to a Format constant.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "raw":
		return FormatRaw, nil
	case "json":
		return FormatJSON, nil
	case "tsv":
		return FormatTSV, nil
	default:
		return FormatRaw, fmt.Errorf("output: unknown format %q (valid: raw, json, tsv)", s)
	}
}

// FormatLine formats a log line according to the chosen Format.
// lineNum is 1-based and used only by formats that include it.
func FormatLine(line string, lineNum int64, f Format) string {
	switch f {
	case FormatJSON:
		escaped := strings.ReplaceAll(line, `"`, `\"`)
		return fmt.Sprintf(`{"n":%d,"line":"%s"}`, lineNum, escaped)
	case FormatTSV:
		return fmt.Sprintf("%d\t%s", lineNum, line)
	default:
		return line
	}
}
