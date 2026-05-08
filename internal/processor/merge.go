package processor

import (
	"io"
	"sort"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/reader"
)

// MergeEntry holds a parsed log line along with its source file name.
type MergeEntry struct {
	Line   parser.LogLine
	Source string
}

// MergeFiles reads all provided log files, parses each line, and writes
// them to the writer sorted by timestamp (ascending). Lines without a
// timestamp are appended after sorted lines in file order.
func MergeFiles(paths []string, w *output.Writer, format string) error {
	p := parser.NewParser()
	var sorted []MergeEntry
	var unsorted []MergeEntry

	for _, path := range paths {
		lr, err := reader.NewLineReader(path)
		if err != nil {
			return err
		}
		for {
			raw, err := lr.ReadLine()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			line := p.Parse(raw)
			entry := MergeEntry{Line: line, Source: path}
			if line.Timestamp.IsZero() {
				unsorted = append(unsorted, entry)
			} else {
				sorted = append(sorted, entry)
			}
		}
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Line.Timestamp.Before(sorted[j].Line.Timestamp)
	})

	all := append(sorted, unsorted...)
	for _, entry := range all {
		formatted := output.FormatLine(entry.Line, format)
		if err := w.WriteLine(formatted); err != nil {
			return err
		}
	}
	return nil
}
