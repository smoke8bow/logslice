package processor

import (
	"errors"
	"sync"
	"time"
)

// RateLimiter limits the number of log lines emitted per second.
type RateLimiter struct {
	mu       sync.Mutex
	rate     int
	window   time.Duration
	bucket   int
	lastTick time.Time
}

// NewRateLimiter creates a RateLimiter that allows at most ratePerSec lines
// per second. ratePerSec must be >= 1.
func NewRateLimiter(ratePerSec int) (*RateLimiter, error) {
	if ratePerSec < 1 {
		return nil, errors.New("ratelimit: rate must be at least 1")
	}
	return &RateLimiter{
		rate:     ratePerSec,
		window:   time.Second,
		bucket:   ratePerSec,
		lastTick: time.Now(),
	}, nil
}

// Allow returns true if the line should be passed through, false if it should
// be dropped to stay within the configured rate limit.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastTick)

	// Refill tokens proportional to elapsed time.
	if elapsed >= r.window {
		r.bucket = r.rate
		r.lastTick = now
	} else {
		refill := int(float64(r.rate) * elapsed.Seconds())
		if refill > 0 {
			r.bucket += refill
			if r.bucket > r.rate {
				r.bucket = r.rate
			}
			r.lastTick = now
		}
	}

	if r.bucket <= 0 {
		return false
	}
	r.bucket--
	return true
}

// Reset resets the limiter to a full token bucket.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bucket = r.rate
	r.lastTick = time.Now()
}
