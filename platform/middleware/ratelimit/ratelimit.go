package ratelimit

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimit struct {
	requestsPerMinute int
	mu                sync.Mutex
	requests          map[string][]time.Time
}

func New(requestsPerMinute int) func(http.Handler) http.Handler {
	rl := &RateLimit{
		requestsPerMinute: requestsPerMinute,
		requests:          make(map[string][]time.Time),
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]

			rl.mu.Lock()
			defer rl.mu.Unlock()

			// Clean up old requests
			now := time.Now()
			cutoff := now.Add(-time.Minute)
			var validRequests []time.Time
			for _, t := range rl.requests[ip] {
				if t.After(cutoff) {
					validRequests = append(validRequests, t)
				}
			}
			rl.requests[ip] = validRequests

			// Check if rate limit exceeded
			if len(rl.requests[ip]) >= rl.requestsPerMinute {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add current request
			rl.requests[ip] = append(rl.requests[ip], now)

			next.ServeHTTP(w, r)
		})
	}
}
