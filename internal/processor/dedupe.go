package processor

import (
	"crypto/fnv"
	"sync"
)

// Deduplicator filters out duplicate log lines using a hash set.
type Deduplicator struct {
	mu   sync.Mutex
	seen map[uint64]struct{}
	max  int
}

// NewDeduplicator creates a Deduplicator with a maximum cache size.
// If max <= 0, no limit is enforced.
func NewDeduplicator(max int) *Deduplicator {
	return &Deduplicator{
		seen: make(map[uint64]struct{}),
		max:  max,
	}
}

// IsDuplicate returns true if the line has been seen before.
func (d *Deduplicator) IsDuplicate(line string) bool {
	h := hash(line)
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.seen[h]; ok {
		return true
	}
	if d.max > 0 && len(d.seen) >= d.max {
		// Evict by clearing — simple strategy to bound memory.
		d.seen = make(map[uint64]struct{})
	}
	d.seen[h] = struct{}{}
	return false
}

// Reset clears the seen set.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[uint64]struct{})
}

// Size returns the number of unique hashes stored.
func (d *Deduplicator) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
