package processor

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewDeduplicator_InitialSize(t *testing.T) {
	d := NewDeduplicator(100)
	if d.Size() != 0 {
		t.Fatalf("expected size 0, got %d", d.Size())
	}
}

func TestDeduplicator_NoDuplicates(t *testing.T) {
	d := NewDeduplicator(100)
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Errorf("line %q should not be a duplicate on first seen", l)
		}
	}
	if d.Size() != 3 {
		t.Fatalf("expected size 3, got %d", d.Size())
	}
}

func TestDeduplicator_DetectsDuplicates(t *testing.T) {
	d := NewDeduplicator(100)
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Error("expected 'hello' to be detected as duplicate")
	}
}

func TestDeduplicator_Reset(t *testing.T) {
	d := NewDeduplicator(100)
	d.IsDuplicate("line1")
	d.IsDuplicate("line2")
	d.Reset()
	if d.Size() != 0 {
		t.Fatalf("expected size 0 after reset, got %d", d.Size())
	}
	if d.IsDuplicate("line1") {
		t.Error("expected 'line1' to not be a duplicate after reset")
	}
}

func TestDeduplicator_MaxEviction(t *testing.T) {
	max := 5
	d := NewDeduplicator(max)
	for i := 0; i < max; i++ {
		d.IsDuplicate(fmt.Sprintf("line-%d", i))
	}
	// Next insert should trigger eviction.
	d.IsDuplicate("overflow")
	if d.Size() > max {
		t.Errorf("size %d exceeds max %d after eviction", d.Size(), max)
	}
}

func TestDeduplicator_ConcurrentSafe(t *testing.T) {
	d := NewDeduplicator(0)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			line := fmt.Sprintf("concurrent-line-%d", n%10)
			d.IsDuplicate(line)
		}(i)
	}
	wg.Wait()
}

func TestDeduplicator_UnlimitedMax(t *testing.T) {
	d := NewDeduplicator(0)
	for i := 0; i < 200; i++ {
		d.IsDuplicate(fmt.Sprintf("entry-%d", i))
	}
	if d.Size() != 200 {
		t.Fatalf("expected 200 unique entries, got %d", d.Size())
	}
}
