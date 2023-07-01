package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type IPRateLimiter struct {
	limiterMap map[string]*rate.Limiter
	mu         sync.Mutex
	rateLimit  rate.Limit
	burstLimit int
}

type RateLimitOptions struct {
	RateLimit  rate.Limit
	BurstLimit int
}

func NewIPRateLimiter(rateLimit rate.Limit, burstLimit int) *IPRateLimiter {
	return &IPRateLimiter{
		limiterMap: make(map[string]*rate.Limiter),
		rateLimit:  rateLimit,
		burstLimit: burstLimit,
	}
}

func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, ok := rl.limiterMap[ip]
	if !ok {
		limiter = rate.NewLimiter(rl.rateLimit, rl.burstLimit)
		rl.limiterMap[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(options RateLimitOptions) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		rl := NewIPRateLimiter(options.RateLimit, options.BurstLimit)
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := rl.GetLimiter(ip)

			if !limiter.Allow() {
				RespondWithJSON(w, 429, "too many requests", nil)
				return
			}

			next(w, r)
		}
	}
}
