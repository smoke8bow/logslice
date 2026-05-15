package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTmplTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "tmpllog-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestRunTemplate_BasicRaw(t *testing.T) {
	input := writeTmplTempLog(t, []string{"hello", "world"})
	out := filepath.Join(t.TempDir(), "out.log")
	err := RunTemplate(TemplateConfig{
		InputFile:  input,
		OutputFile: out,
		Template:   "LINE: {{.Raw}}",
		Format:     "raw",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "LINE: hello") {
		t.Errorf("expected 'LINE: hello' in output, got: %s", data)
	}
}

func TestRunTemplate_MissingFile(t *testing.T) {
	err := RunTemplate(TemplateConfig{
		InputFile: "/no/such/file.log",
		Template:  "{{.Raw}}",
		Format:    "raw",
	})
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestRunTemplate_InvalidTemplate(t *testing.T) {
	input := writeTmplTempLog(t, []string{"line"})
	err := RunTemplate(TemplateConfig{
		InputFile: input,
		Template:  "",
		Format:    "raw",
	})
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestRunTemplate_InvalidOutputPath(t *testing.T) {
	input := writeTmplTempLog(t, []string{"line"})
	err := RunTemplate(TemplateConfig{
		InputFile:  input,
		OutputFile: "/no/such/dir/out.log",
		Template:   "{{.Raw}}",
		Format:     "raw",
	})
	if err == nil {
		t.Fatal("expected error for invalid output path")
	}
}

func TestRunTemplate_InvalidFormat(t *testing.T) {
	input := writeTmplTempLog(t, []string{"line"})
	out := filepath.Join(t.TempDir(), "out.log")
	err := RunTemplate(TemplateConfig{
		InputFile:  input,
		OutputFile: out,
		Template:   "{{.Raw}}",
		Format:     "badformat",
	})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
