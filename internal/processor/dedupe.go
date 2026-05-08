package processor

import (
	"hash/fnv"
	"sync"
)

// Deduplicator tracks seen log lines and filters out duplicates within a
// sliding window of recently observed content hashes. It is safe for
// concurrent use by multiple goroutines.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[uint64]struct{}
	maxSize int
}

// NewDeduplicator creates a Deduplicator that remembers up to maxSize unique
// lines. When the internal set exceeds maxSize it is cleared, providing a
// simple bounded-memory rolling-window deduplication strategy.
//
// A maxSize of 0 disables deduplication (IsDuplicate always returns false).
func NewDeduplicator(maxSize int) *Deduplicator {
	return &Deduplicator{
		seen:    make(map[uint64]struct{}, maxSize),
		maxSize: maxSize,
	}
}

// IsDuplicate reports whether line has been seen before. If it has not been
// seen, it is recorded and false is returned. If maxSize is 0 the method
// always returns false without recording anything.
func (d *Deduplicator) IsDuplicate(line string) bool {
	if d.maxSize == 0 {
		return false
	}

	h := hash(line)

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.seen[h]; ok {
		return true
	}

	// Evict the whole window when the cap is reached to keep memory bounded.
	if len(d.seen) >= d.maxSize {
		d.seen = make(map[uint64]struct{}, d.maxSize)
	}

	d.seen[h] = struct{}{}
	return false
}

// Reset clears all recorded hashes, allowing previously seen lines to pass
// through again.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[uint64]struct{}, d.maxSize)
}

// Len returns the number of unique hashes currently stored.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// hash returns a 64-bit FNV-1a hash of s.
func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
