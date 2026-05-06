package parser

import (
	"encoding/json"
	"time"
)

var commonTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
}

func parseJSON(raw string, line *LogLine, timeFields []string) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		line.Message = raw
		return
	}
	for k, v := range m {
		s := toString(v)
		line.Fields[k] = s
	}
	// Extract timestamp
	for _, tf := range timeFields {
		if v, ok := m[tf]; ok {
			if t, err := parseTime(toString(v)); err == nil {
				line.Timestamp = t
				break
			}
		}
	}
	// Extract level
	for _, lf := range []string{"level", "lvl", "severity"} {
		if v, ok := m[lf]; ok {
			line.Level = toString(v)
			break
		}
	}
	// Extract message
	for _, mf := range []string{"msg", "message", "text"} {
		if v, ok := m[mf]; ok {
			line.Message = toString(v)
			break
		}
	}
}

func parseTime(s string) (time.Time, error) {
	for _, layout := range commonTimeLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errNoMatch
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case nil:
		return ""
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}
