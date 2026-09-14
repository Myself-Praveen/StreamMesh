package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket implements a simple thread-safe token bucket rate limiter
type TokenBucket struct {
	capacity  int
	tokens    int
	rate      float64 // tokens per second
	lastCheck time.Time
	mu        sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:  capacity,
		tokens:    capacity,
		rate:      rate,
		lastCheck: time.Now(),
	}
}

// Allow checks if a request is allowed. Returns true if allowed.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastCheck).Seconds()

	// Replenish tokens
	tb.tokens += int(elapsed * tb.rate)
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastCheck = now

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
