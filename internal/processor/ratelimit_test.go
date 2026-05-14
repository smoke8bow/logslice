package processor

import (
	"testing"
	"time"
)

func TestNewRateLimiter_InvalidRate(t *testing.T) {
	_, err := NewRateLimiter(0)
	if err == nil {
		t.Fatal("expected error for rate=0, got nil")
	}
	_, err = NewRateLimiter(-5)
	if err == nil {
		t.Fatal("expected error for rate=-5, got nil")
	}
}

func TestNewRateLimiter_ValidRate(t *testing.T) {
	rl, err := NewRateLimiter(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rl == nil {
		t.Fatal("expected non-nil RateLimiter")
	}
}

func TestRateLimiter_AllowsUpToRate(t *testing.T) {
	const rate = 5
	rl, _ := NewRateLimiter(rate)

	allowed := 0
	for i := 0; i < rate*2; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	// Should allow exactly `rate` lines before bucket is empty.
	if allowed != rate {
		t.Errorf("expected %d allowed, got %d", rate, allowed)
	}
}

func TestRateLimiter_DropsWhenExhausted(t *testing.T) {
	rl, _ := NewRateLimiter(2)

	// Drain the bucket.
	rl.Allow()
	rl.Allow()

	if rl.Allow() {
		t.Error("expected Allow() to return false when bucket is empty")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl, _ := NewRateLimiter(3)

	// Drain completely.
	rl.Allow()
	rl.Allow()
	rl.Allow()

	if rl.Allow() {
		t.Fatal("bucket should be empty before reset")
	}

	rl.Reset()

	if !rl.Allow() {
		t.Error("expected Allow() to return true after Reset")
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	rl, _ := NewRateLimiter(100)

	// Drain the full bucket.
	for i := 0; i < 100; i++ {
		rl.Allow()
	}
	if rl.Allow() {
		t.Fatal("bucket should be empty after draining")
	}

	// Wait for at least a partial refill (> 50ms = ~5 tokens at 100/s).
	time.Sleep(60 * time.Millisecond)

	if !rl.Allow() {
		t.Error("expected at least one token to be refilled after 60ms")
	}
}
