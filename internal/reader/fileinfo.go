package reader

import (
	"fmt"
	"os"
)

// FileInfo holds metadata about a log file relevant to slicing.
type FileInfo struct {
	Path string
	Size int64
}

// StatFile returns a FileInfo for the given path.
func StatFile(path string) (FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, fmt.Errorf("stat %q: %w", path, err)
	}
	if fi.IsDir() {
		return FileInfo{}, fmt.Errorf("%q is a directory, not a file", path)
	}
	return FileInfo{Path: path, Size: fi.Size()}, nil
}

// ChunkOffsets divides a file into n roughly equal byte-offset chunks.
// Each returned value is a start offset; callers read until the next offset
// (or EOF for the last chunk). Returns at least one offset (0).
func ChunkOffsets(size int64, n int) []int64 {
	if n <= 1 || size == 0 {
		return []int64{0}
	}
	chunkSize := size / int64(n)
	offsets := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		off := int64(i) * chunkSize
		if off >= size {
			break
		}
		offsets = append(offsets, off)
	}
	return offsets
}
