package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatFile_Valid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "app.log")
	if err := os.WriteFile(p, []byte("hello world"), 0o644); err != nil {
		t.Fatal(err)
	}
	fi, err := StatFile(p)
	if err != nil {
		t.Fatalf("StatFile: %v", err)
	}
	if fi.Path != p {
		t.Errorf("expected path %q, got %q", p, fi.Path)
	}
	if fi.Size != 11 {
		t.Errorf("expected size 11, got %d", fi.Size)
	}
}

func TestStatFile_Missing(t *testing.T) {
	_, err := StatFile("/no/such/file.log")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestStatFile_Directory(t *testing.T) {
	dir := t.TempDir()
	_, err := StatFile(dir)
	if err == nil {
		t.Fatal("expected error for directory path")
	}
}

func TestChunkOffsets_SingleChunk(t *testing.T) {
	offsets := ChunkOffsets(1000, 1)
	if len(offsets) != 1 || offsets[0] != 0 {
		t.Errorf("expected [0], got %v", offsets)
	}
}

func TestChunkOffsets_MultipleChunks(t *testing.T) {
	offsets := ChunkOffsets(1000, 4)
	if len(offsets) != 4 {
		t.Fatalf("expected 4 offsets, got %d", len(offsets))
	}
	if offsets[0] != 0 {
		t.Errorf("first offset should be 0, got %d", offsets[0])
	}
	if offsets[1] != 250 {
		t.Errorf("second offset should be 250, got %d", offsets[1])
	}
}

func TestChunkOffsets_ZeroSize(t *testing.T) {
	offsets := ChunkOffsets(0, 4)
	if len(offsets) != 1 || offsets[0] != 0 {
		t.Errorf("expected [0] for zero-size file, got %v", offsets)
	}
}

func TestChunkOffsets_MoreChunksThanBytes(t *testing.T) {
	offsets := ChunkOffsets(3, 10)
	if len(offsets) > 3 {
		t.Errorf("expected at most 3 offsets for 3-byte file, got %d", len(offsets))
	}
}
