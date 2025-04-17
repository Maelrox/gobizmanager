package ratelimiter

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	attempts    map[string][]time.Time
	bannedIPs   map[string]time.Time
	banTime     time.Duration
	maxAttempts int
	windowSize  time.Duration
}

func NewRateLimiter(maxAttempts int, windowSize, banTime time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts:    make(map[string][]time.Time),
		bannedIPs:   make(map[string]time.Time),
		banTime:     banTime,
		maxAttempts: maxAttempts,
		windowSize:  windowSize,
	}
}

func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if banTime, banned := rl.bannedIPs[ip]; banned {
		if now.Sub(banTime) < rl.banTime {
			return false
		}
		// Ban period has expired, remove from banned list
		delete(rl.bannedIPs, ip)
	}

	attempts, exists := rl.attempts[ip]

	// Clean up old attempts
	if exists {
		var validAttempts []time.Time
		for _, attempt := range attempts {
			if now.Sub(attempt) <= rl.windowSize {
				validAttempts = append(validAttempts, attempt)
			}
		}
		rl.attempts[ip] = validAttempts
		attempts = validAttempts
	}

	// Check if IP has exceeded max attempts
	if len(attempts) >= rl.maxAttempts {
		rl.bannedIPs[ip] = now
		delete(rl.attempts, ip)
		return false
	}

	// Record new attempt
	rl.attempts[ip] = append(attempts, now)
	return true
}

func (rl *RateLimiter) Reset(ip string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.attempts, ip)
	delete(rl.bannedIPs, ip)
}
