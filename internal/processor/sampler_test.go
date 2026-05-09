package processor

import (
	"testing"
)

func TestNewSampler_InvalidRate(t *testing.T) {
	_, err := NewSampler("rate", 0.0, 0)
	if err == nil {
		t.Fatal("expected error for rate=0.0")
	}
	_, err = NewSampler("rate", 1.5, 0)
	if err == nil {
		t.Fatal("expected error for rate=1.5")
	}
}

func TestNewSampler_InvalidNth(t *testing.T) {
	_, err := NewSampler("nth", 0, 0)
	if err == nil {
		t.Fatal("expected error for nth=0")
	}
}

func TestNewSampler_UnknownMode(t *testing.T) {
	_, err := NewSampler("random", 0.5, 1)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestSampler_NthMode_KeepsEveryNth(t *testing.T) {
	s, err := NewSampler("nth", 0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var kept []int
	for i := 1; i <= 9; i++ {
		if s.Keep() {
			kept = append(kept, i)
		}
	}
	expected := []int{3, 6, 9}
	if len(kept) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, kept)
	}
	for i, v := range expected {
		if kept[i] != v {
			t.Errorf("position %d: expected %d, got %d", i, v, kept[i])
		}
	}
}

func TestSampler_Reset_ResetsCounter(t *testing.T) {
	s, _ := NewSampler("nth", 0, 2)
	s.Keep() // 1 — skip
	s.Keep() // 2 — keep
	s.Reset()
	s.Keep() // 1 again — skip
	if s.Keep() { // 2 again — keep
		// expected
	} else {
		t.Error("expected Keep() to return true after Reset at count=2")
	}
}

func TestSampler_RateMode_ApproximateSampling(t *testing.T) {
	s, err := NewSampler("rate", 0.5, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const total = 10000
	kept := 0
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	// Allow 10% tolerance around 50%
	if kept < total*40/100 || kept > total*60/100 {
		t.Errorf("rate=0.5 kept %d/%d lines, outside expected range", kept, total)
	}
}

func TestSampler_RateMode_FullRate(t *testing.T) {
	s, err := NewSampler("rate", 1.0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatalf("rate=1.0 should keep all lines, failed at line %d", i+1)
		}
	}
}
