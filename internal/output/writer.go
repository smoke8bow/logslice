package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Writer wraps an io.Writer with buffered output and optional line counting.
type Writer struct {
	bw        *bufio.Writer
	lineCount int64
	closer    io.Closer
}

// NewWriter creates a Writer that writes to the given io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		bw: bufio.NewWriterSize(w, 64*1024),
	}
}

// NewFileWriter opens (or creates) a file at path and returns a Writer for it.
func NewFileWriter(path string) (*Writer, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: cannot open file %q: %w", path, err)
	}
	return &Writer{
		bw:     bufio.NewWriterSize(f, 64*1024),
		closer: f,
	}, nil
}

// WriteLine writes a single line followed by a newline character.
func (w *Writer) WriteLine(line string) error {
	if _, err := fmt.Fprintln(w.bw, line); err != nil {
		return fmt.Errorf("output: write error: %w", err)
	}
	w.lineCount++
	return nil
}

// WriteBytes writes raw bytes followed by a newline.
func (w *Writer) WriteBytes(b []byte) error {
	if _, err := w.bw.Write(b); err != nil {
		return fmt.Errorf("output: write error: %w", err)
	}
	if err := w.bw.WriteByte('\n'); err != nil {
		return fmt.Errorf("output: write newline error: %w", err)
	}
	w.lineCount++
	return nil
}

// LineCount returns the number of lines written so far.
func (w *Writer) LineCount() int64 {
	return w.lineCount
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	if err := w.bw.Flush(); err != nil {
		return fmt.Errorf("output: flush error: %w", err)
	}
	return nil
}

// Close flushes and closes the underlying writer if it is an io.Closer.
func (w *Writer) Close() error {
	if err := w.Flush(); err != nil {
		return err
	}
	if w.closer != nil {
		return w.closer.Close()
	}
	return nil
}
