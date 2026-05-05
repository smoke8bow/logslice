package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// PatternFilter holds compiled patterns for log line matching.
type PatternFilter struct {
	includes []*regexp.Regexp
	excludes []*regexp.Regexp
}

// NewPatternFilter creates a PatternFilter from include and exclude pattern strings.
// Patterns are treated as regular expressions.
func NewPatternFilter(includes, excludes []string) (*PatternFilter, error) {
	pf := &PatternFilter{}

	for _, p := range includes {
		if strings.TrimSpace(p) == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid include pattern %q: %w", p, err)
		}
		pf.includes = append(pf.includes, re)
	}

	for _, p := range excludes {
		if strings.TrimSpace(p) == "" {
			continue
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern %q: %w", p, err)
		}
		pf.excludes = append(pf.excludes, re)
	}

	return pf, nil
}

// Match returns true if the line passes the pattern filter.
// A line passes if it matches at least one include pattern (or no includes are set)
// and does not match any exclude pattern.
func (pf *PatternFilter) Match(line string) bool {
	if len(pf.excludes) > 0 {
		for _, re := range pf.excludes {
			if re.MatchString(line) {
				return false
			}
		}
	}

	if len(pf.includes) == 0 {
		return true
	}

	for _, re := range pf.includes {
		if re.MatchString(line) {
			return true
		}
	}

	return false
}

// HasFilters reports whether any include or exclude patterns are configured.
func (pf *PatternFilter) HasFilters() bool {
	return len(pf.includes) > 0 || len(pf.excludes) > 0
}
