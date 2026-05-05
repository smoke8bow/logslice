package reader

import (
	"bufio"
	"io"
	"os"
)

// LineReader reads lines from a log file with optional byte-range seeking.
type LineReader struct {
	file    *os.File
	scanner *bufio.Scanner
	path    string
}

// NewLineReader opens the file at path and returns a LineReader.
func NewLineReader(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{file: f, scanner: scanner, path: path}, nil
}

// NewLineReaderAt opens the file at path and seeks to the given byte offset
// before scanning. Useful for splitting large files across goroutines.
func NewLineReaderAt(path string, offset int64) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			f.Close()
			return nil, err
		}
	}
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{file: f, scanner: scanner, path: path}, nil
}

// Scan advances to the next line. Returns false when done or on error.
func (r *LineReader) Scan() bool {
	return r.scanner.Scan()
}

// Text returns the current line text.
func (r *LineReader) Text() string {
	return r.scanner.Text()
}

// Err returns any scanning error (excluding io.EOF).
func (r *LineReader) Err() error {
	return r.scanner.Err()
}

// Close releases the underlying file handle.
func (r *LineReader) Close() error {
	return r.file.Close()
}

// Path returns the file path this reader was opened with.
func (r *LineReader) Path() string {
	return r.path
}
