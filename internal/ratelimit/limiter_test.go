package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	// 2 tokens burst, 10 tokens per second
	tb := NewTokenBucket(2, 10.0)

	// Consume burst
	if !tb.Allow() {
		t.Errorf("Expected first request to be allowed")
	}
	if !tb.Allow() {
		t.Errorf("Expected second request to be allowed")
	}

	// Third request should be denied immediately
	if tb.Allow() {
		t.Errorf("Expected third request to be denied")
	}

	// Wait for 1 token to regenerate (0.1s)
	time.Sleep(110 * time.Millisecond)

	if !tb.Allow() {
		t.Errorf("Expected request to be allowed after replenish")
	}
}

func TestIPRateLimiter(t *testing.T) {
	limiter := NewIPRateLimiter(1, 1.0)
	ip := "127.0.0.1"

	if !limiter.Allow(ip) {
		t.Errorf("Expected first request to be allowed")
	}

	if limiter.Allow(ip) {
		t.Errorf("Expected second request from same IP to be denied")
	}

	if !limiter.Allow("192.168.1.1") {
		t.Errorf("Expected request from different IP to be allowed")
	}
}
