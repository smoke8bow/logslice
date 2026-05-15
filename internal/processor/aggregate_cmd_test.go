package processor

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func writeAggTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "aggtest-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()
	f.WriteString(strings.Join(lines, "\n") + "\n")
	return f.Name()
}

func TestRunAggregate_CountRaw(t *testing.T) {
	lines := []string{
		`level=info msg="request handled"`,
		`level=info msg="another"`,
		`level=error msg="oops"`,
	}
	path := writeAggTempLog(t, lines)
	defer os.Remove(path)

	var buf bytes.Buffer
	err := RunAggregate(AggregateConfig{
		InputFile:  path,
		GroupField: "level",
		Mode:       AggCount,
		OutputFmt:  output.FormatRaw,
		Writer:     &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "info") || !strings.Contains(out, "error") {
		t.Errorf("expected both groups in output, got: %s", out)
	}
}

func TestRunAggregate_SumJSON(t *testing.T) {
	lines := []string{
		`{"service":"api","latency":"12"}`,
		`{"service":"api","latency":"8"}`,
		`{"service":"db","latency":"5"}`,
	}
	path := writeAggTempLog(t, lines)
	defer os.Remove(path)

	var buf bytes.Buffer
	err := RunAggregate(AggregateConfig{
		InputFile:  path,
		GroupField: "service",
		ValueField: "latency",
		Mode:       AggSum,
		OutputFmt:  output.FormatJSON,
		Writer:     &buf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"api"`) {
		t.Errorf("expected api in output, got: %s", out)
	}
	if !strings.Contains(out, "20") {
		t.Errorf("expected sum 20 in output, got: %s", out)
	}
}

func TestRunAggregate_MissingFile(t *testing.T) {
	var buf bytes.Buffer
	err := RunAggregate(AggregateConfig{
		InputFile:  "/no/such/file.log",
		GroupField: "level",
		Mode:       AggCount,
		OutputFmt:  output.FormatRaw,
		Writer:     &buf,
	})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunAggregate_InvalidAggregator(t *testing.T) {
	var buf bytes.Buffer
	err := RunAggregate(AggregateConfig{
		InputFile:  "unused",
		GroupField: "",
		Mode:       AggCount,
		OutputFmt:  output.FormatRaw,
		Writer:     &buf,
	})
	if err == nil {
		t.Fatal("expected error for empty groupField")
	}
}
