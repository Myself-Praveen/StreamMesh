package ratelimit

import (
	"net"
	"net/http"
	"sync"
)

// IPRateLimiter limits connections based on IP address
type IPRateLimiter struct {
	limiters map[string]*TokenBucket
	mu       sync.RWMutex
	burst    int
	rate     float64
}

// NewIPRateLimiter creates a new IPRateLimiter
func NewIPRateLimiter(burst int, rate float64) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*TokenBucket),
		burst:    burst,
		rate:     rate,
	}
}

// getLimiter returns the limiter for an IP
func (i *IPRateLimiter) getLimiter(ip string) *TokenBucket {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double check
	limiter, exists = i.limiters[ip]
	if exists {
		return limiter
	}

	limiter = NewTokenBucket(i.burst, i.rate)
	i.limiters[ip] = limiter
	return limiter
}

// Allow checks if the request from the IP is allowed
func (i *IPRateLimiter) Allow(ip string) bool {
	return i.getLimiter(ip).Allow()
}

// GetIP extracts the IP address from a request
func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}
