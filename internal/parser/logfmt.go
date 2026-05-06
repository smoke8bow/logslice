package parser

import (
	"errors"
	"strings"
)

var errNoMatch = errors.New("no match")

// parseLogfmt parses a logfmt-style line (key=value key="value with spaces").
func parseLogfmt(raw string, line *LogLine, timeFields []string) {
	fields := splitLogfmt(raw)
	for k, v := range fields {
		line.Fields[k] = v
	}
	// Extract timestamp
	for _, tf := range timeFields {
		if v, ok := fields[tf]; ok {
			if t, err := parseTime(v); err == nil {
				line.Timestamp = t
				break
			}
		}
	}
	// Extract level
	for _, lf := range []string{"level", "lvl", "severity"} {
		if v, ok := fields[lf]; ok {
			line.Level = v
			break
		}
	}
	// Extract message
	for _, mf := range []string{"msg", "message", "text"} {
		if v, ok := fields[mf]; ok {
			line.Message = v
			break
		}
	}
	if line.Message == "" {
		line.Message = raw
	}
}

// splitLogfmt splits a logfmt line into key/value pairs.
func splitLogfmt(raw string) map[string]string {
	result := make(map[string]string)
	s := strings.TrimSpace(raw)
	for len(s) > 0 {
		// Find key
		eqIdx := strings.IndexByte(s, '=')
		if eqIdx < 0 {
			break
		}
		key := strings.TrimSpace(s[:eqIdx])
		s = s[eqIdx+1:]
		var value string
		if len(s) > 0 && s[0] == '"' {
			// Quoted value
			end := strings.Index(s[1:], "\"")
			if end < 0 {
				value = s[1:]
				s = ""
			} else {
				value = s[1 : end+1]
				s = strings.TrimSpace(s[end+2:])
			}
		} else {
			spIdx := strings.IndexByte(s, ' ')
			if spIdx < 0 {
				value = s
				s = ""
			} else {
				value = s[:spIdx]
				s = strings.TrimSpace(s[spIdx+1:])
			}
		}
		if key != "" {
			result[key] = value
		}
	}
	return result
}
