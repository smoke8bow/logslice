package reader

import (
	"bufio"
	"io"
	"time"
)

// LogLine represents a single parsed log line with metadata.
type LogLine struct {
	Raw       string
	Timestamp time.Time
	Offset    int64
	LineNum   int
}

// Scanner wraps a LineReader and emits LogLine values.
type Scanner struct {
	reader  io.Reader
	scanner *bufio.Scanner
	lineNum int
	offset  int64
}

// NewScanner creates a Scanner from an io.Reader.
func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &Scanner{
		reader:  r,
		scanner: s,
	}
}

// Scan advances to the next line. Returns false when done or on error.
func (s *Scanner) Scan() bool {
	return s.scanner.Scan()
}

// Line returns the current LogLine.
func (s *Scanner) Line() LogLine {
	raw := s.scanner.Text()
	s.lineNum++
	prev := s.offset
	s.offset += int64(len(raw)) + 1 // +1 for newline
	return LogLine{
		Raw:     raw,
		Offset:  prev,
		LineNum: s.lineNum,
	}
}

// Err returns any error encountered during scanning.
func (s *Scanner) Err() error {
	return s.scanner.Err()
}
