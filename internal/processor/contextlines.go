package processor

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// ContextLines holds lines before and after a matching line.
type ContextLines struct {
	before int
	after  int
	buf    []parser.LogLine
	pending []parser.LogLine
	pendingCount int
}

// NewContextLines creates a ContextLines processor that emits `before` lines
// before a match and `after` lines after a match. A line is considered a
// "match" when the supplied match function returns true.
func NewContextLines(before, after int) (*ContextLines, error) {
	if before < 0 {
		return nil, fmt.Errorf("contextlines: before must be >= 0, got %d", before)
	}
	if after < 0 {
		return nil, fmt.Errorf("contextlines: after must be >= 0, got %d", after)
	}
	return &ContextLines{
		before: before,
		after:  after,
		buf:    make([]parser.LogLine, 0, before),
	}, nil
}

// Process accepts a line and a matched flag. It returns the set of lines that
// should be emitted as a result (may be empty, or contain context lines).
func (c *ContextLines) Process(line parser.LogLine, matched bool) []parser.LogLine {
	var out []parser.LogLine

	if matched {
		// Flush buffered before-context
		out = append(out, c.buf...)
		c.buf = c.buf[:0]
		out = append(out, line)
		c.pendingCount = c.after
		c.pending = c.pending[:0]
		return out
	}

	if c.pendingCount > 0 {
		out = append(out, line)
		c.pendingCount--
		return out
	}

	// Buffer for before-context
	if c.before > 0 {
		if len(c.buf) >= c.before {
			c.buf = c.buf[1:]
		}
		c.buf = append(c.buf, line)
	}

	return out
}

// Reset clears internal state between processing runs.
func (c *ContextLines) Reset() {
	c.buf = c.buf[:0]
	c.pending = c.pending[:0]
	c.pendingCount = 0
}
