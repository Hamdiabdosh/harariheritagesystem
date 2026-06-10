package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter(maxTokens, refillRate float64) *gin.Engine {
	r := gin.New()
	r.POST("/auth/login", RateLimit(maxTokens, refillRate), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func doRequest(r *gin.Engine, ip string) int {
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set("X-Forwarded-For", ip)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w.Code
}

// TestRateLimit_AllowsBurst verifies that maxTokens requests succeed immediately.
func TestRateLimit_AllowsBurst(t *testing.T) {
	r := newTestRouter(5, 0.1)
	ip := "10.0.0.1"

	for i := range 5 {
		code := doRequest(r, ip)
		if code != http.StatusOK {
			t.Fatalf("request %d: want 200, got %d", i+1, code)
		}
	}
}

// TestRateLimit_BlocksAfterBurst verifies the (maxTokens+1)th request is rejected.
func TestRateLimit_BlocksAfterBurst(t *testing.T) {
	r := newTestRouter(3, 0.1)
	ip := "10.0.0.2"

	for range 3 {
		doRequest(r, ip)
	}

	code := doRequest(r, ip)
	if code != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", code)
	}
}

// TestRateLimit_DifferentIPsAreIndependent verifies two IPs don't share a bucket.
func TestRateLimit_DifferentIPsAreIndependent(t *testing.T) {
	r := newTestRouter(2, 0.1)

	// Exhaust IP A
	doRequest(r, "10.0.0.10")
	doRequest(r, "10.0.0.10")
	if code := doRequest(r, "10.0.0.10"); code != http.StatusTooManyRequests {
		t.Fatalf("IP A should be rate limited, got %d", code)
	}

	// IP B should still be free
	if code := doRequest(r, "10.0.0.11"); code != http.StatusOK {
		t.Fatalf("IP B should not be rate limited, got %d", code)
	}
}

// TestRateLimit_RefillsOverTime verifies tokens come back after waiting.
func TestRateLimit_RefillsOverTime(t *testing.T) {
	// refillRate = 100 tokens/sec so we don't need to wait long in tests
	r := newTestRouter(1, 100)
	ip := "10.0.0.3"

	// Use the one token
	if code := doRequest(r, ip); code != http.StatusOK {
		t.Fatalf("first request: want 200, got %d", code)
	}
	// Immediately blocked
	if code := doRequest(r, ip); code != http.StatusTooManyRequests {
		t.Fatalf("second request: want 429, got %d", code)
	}

	// Wait 20ms — at 100 tokens/sec that's 2 tokens refilled
	time.Sleep(20 * time.Millisecond)

	if code := doRequest(r, ip); code != http.StatusOK {
		t.Fatalf("after refill: want 200, got %d", code)
	}
}

// TestRateLimit_RetryAfterHeader verifies the header is set on 429 responses.
func TestRateLimit_RetryAfterHeader(t *testing.T) {
	r := newTestRouter(1, 0.1)
	ip := "10.0.0.4"

	doRequest(r, ip) // consume the 1 token

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set("X-Forwarded-For", ip)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("want 429, got %d", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Fatal("want Retry-After header on 429 response")
	}
}
