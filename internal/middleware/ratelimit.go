package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ipBucket tracks the state of one IP address in the token bucket.
type ipBucket struct {
	tokens   float64   // current token count (float so partial tokens accumulate)
	lastSeen time.Time // when this bucket was last touched (for cleanup)
	mu       sync.Mutex
}

// rateLimiter holds all per-IP buckets and the shared configuration.
type rateLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*ipBucket
	maxTokens  float64       // burst ceiling (e.g. 10 attempts)
	refillRate float64       // tokens added per second (e.g. 0.1 = 1 per 10s)
	cleanupTTL time.Duration // evict buckets unseen for this long
}

func newRateLimiter(maxTokens float64, refillRate float64, cleanupTTL time.Duration) *rateLimiter {
	rl := &rateLimiter{
		buckets:    make(map[string]*ipBucket),
		maxTokens:  maxTokens,
		refillRate: refillRate,
		cleanupTTL: cleanupTTL,
	}
	// Background goroutine evicts stale buckets every 5 minutes to prevent
	// unbounded memory growth in long-running deployments.
	go rl.cleanup()
	return rl
}

// allow returns true if the request from ip should be allowed through.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	bucket, exists := rl.buckets[ip]
	if !exists {
		bucket = &ipBucket{
			tokens:   rl.maxTokens,
			lastSeen: time.Now(),
		}
		rl.buckets[ip] = bucket
	}
	rl.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastSeen).Seconds()
	bucket.lastSeen = now

	// Refill tokens based on time elapsed since last request.
	bucket.tokens += elapsed * rl.refillRate
	if bucket.tokens > rl.maxTokens {
		bucket.tokens = rl.maxTokens
	}

	if bucket.tokens < 1 {
		return false // rate limit exceeded
	}

	bucket.tokens--
	return true
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, bucket := range rl.buckets {
			bucket.mu.Lock()
			stale := time.Since(bucket.lastSeen) > rl.cleanupTTL
			bucket.mu.Unlock()
			if stale {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit returns a Gin middleware that enforces a token-bucket rate limit
// per client IP address.
//
// Parameters:
//   - maxTokens:  burst ceiling — how many requests an IP can fire at once
//                 before being throttled. Set to 10 for auth endpoints.
//   - refillRate: tokens restored per second. 0.1 = 1 token per 10 seconds,
//                 meaning a sustained rate of 6 requests per minute.
//
// For auth endpoints (login + refresh) we use:
//   - maxTokens  = 10   (allow a short burst, e.g. user retrying quickly)
//   - refillRate = 0.1  (sustained: 1 attempt per 10 seconds per IP)
//
// A legitimate user who misremembers their password can retry ~10 times
// immediately, then once every 10 seconds. An attacker runs out in seconds.
func RateLimit(maxTokens float64, refillRate float64) gin.HandlerFunc {
	rl := newRateLimiter(maxTokens, refillRate, 1*time.Hour)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.allow(ip) {
			c.Header("Retry-After", "10")
			RespondError(c, http.StatusTooManyRequests, "Too many requests — please wait before trying again")
			return
		}

		c.Next()
	}
}
