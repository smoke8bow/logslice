package processor

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// HeadFilter passes only the first N lines through the pipeline.
type HeadFilter struct {
	max   int
	count int
}

// NewHeadFilter creates a HeadFilter that allows at most n lines.
// Returns an error if n is less than 1.
func NewHeadFilter(n int) (*HeadFilter, error) {
	if n < 1 {
		return nil, fmt.Errorf("head: n must be at least 1, got %d", n)
	}
	return &HeadFilter{max: n}, nil
}

// Apply returns true (keep) for the first n lines, false thereafter.
func (h *HeadFilter) Apply(line *parser.LogLine) bool {
	if h.count >= h.max {
		return false
	}
	h.count++
	return true
}

// Reset resets the internal counter so the filter can be reused.
func (h *HeadFilter) Reset() {
	h.count = 0
}

// TailFilter buffers the last N lines and emits them at flush time.
type TailFilter struct {
	max int
	buf []*parser.LogLine
}

// NewTailFilter creates a TailFilter that retains at most n lines.
// Returns an error if n is less than 1.
func NewTailFilter(n int) (*TailFilter, error) {
	if n < 1 {
		return nil, fmt.Errorf("tail: n must be at least 1, got %d", n)
	}
	return &TailFilter{max: n, buf: make([]*parser.LogLine, 0, n)}, nil
}

// Collect adds a line to the rolling buffer.
func (t *TailFilter) Collect(line *parser.LogLine) {
	if len(t.buf) >= t.max {
		t.buf = t.buf[1:]
	}
	t.buf = append(t.buf, line)
}

// Lines returns the buffered tail lines in order.
func (t *TailFilter) Lines() []*parser.LogLine {
	out := make([]*parser.LogLine, len(t.buf))
	copy(out, t.buf)
	return out
}

// Reset clears the internal buffer.
func (t *TailFilter) Reset() {
	t.buf = t.buf[:0]
}
