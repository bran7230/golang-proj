package server

import (
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

func AuthMiddleware(originalHandler http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != apiKey {
			err := SendError(w, http.StatusUnauthorized, "error, invalid auth key")
			if err != nil {
				fmt.Printf("auth middleware error: %s\n", err.Error())
				return
			}
			return
		}
		originalHandler.ServeHTTP(w, r)
	})
}

// RateLimiter tracks rate limiters per API key
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

// GetLimiter returns a rate limiter for a given API key, creating one if it doesn't exist
// rps = requests per second, burst = maximum burst size
func (rl *RateLimiter) GetLimiter(key string, rps float64, burst int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.limiters[key]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	rl.limiters[key] = limiter
	return limiter
}

// RateLimitMiddleware enforces rate limiting per API key
func RateLimitMiddleware(rateLimiter *RateLimiter, rps float64, burst int) func(http.Handler) http.Handler {
	return func(originalHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-key")
			if apiKey == "" {
				// If no API key, use a default limiter for unauthenticated requests
				apiKey = "unauthenticated"
			}

			limiter := rateLimiter.GetLimiter(apiKey, rps, burst)

			if !limiter.Allow() {
				err := SendError(w, http.StatusTooManyRequests, "rate limit exceeded, try again later")
				if err != nil {
					fmt.Printf("rate limit middleware error: %s\n", err.Error())
					return
				}
				return
			}

			originalHandler.ServeHTTP(w, r)
		})
	}
}
