package processor

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/user/logslice/internal/output"
)

func writeMergeLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "merge-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(strings.Join(lines, "\n") + "\n")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestMergeFiles_SortsByTimestamp(t *testing.T) {
	file1 := writeMergeLog(t, []string{
		`{"time":"2024-01-01T10:00:00Z","msg":"first"}`,
		`{"time":"2024-01-01T12:00:00Z","msg":"third"}`,
	})
	file2 := writeMergeLog(t, []string{
		`{"time":"2024-01-01T11:00:00Z","msg":"second"}`,
	})

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	if err := MergeFiles([]string{file1, file2}, w, "raw"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "first") {
		t.Errorf("expected first line to contain 'first', got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "second") {
		t.Errorf("expected second line to contain 'second', got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "third") {
		t.Errorf("expected third line to contain 'third', got: %s", lines[2])
	}
}

func TestMergeFiles_UntimestampedAppendedLast(t *testing.T) {
	file1 := writeMergeLog(t, []string{
		"no timestamp here",
		`{"time":"2024-01-01T09:00:00Z","msg":"timed"}`,
	})

	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	if err := MergeFiles([]string{file1}, w, "raw"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "timed") {
		t.Errorf("expected timed line first, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "no timestamp") {
		t.Errorf("expected untimed line last, got: %s", lines[1])
	}
}

func TestMergeFiles_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	err := MergeFiles([]string{"/nonexistent/file.log"}, w, "raw")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestMergeFiles_EmptyList(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf)
	if err := MergeFiles([]string{}, w, "raw"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}
