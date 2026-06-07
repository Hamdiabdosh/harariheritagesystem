package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Gin registers /records/immovable/:id as a leaf node; suffix GET routes must be
// registered first or /records/immovable/:id/pdf returns 404.
func TestRecordSubResourceRoutes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/records/immovable/:id/pdf", withRecordType("immovable", func(c *gin.Context) {
		c.String(http.StatusOK, "pdf")
	}))
	r.GET("/records/immovable/:id/history", withRecordType("immovable", func(c *gin.Context) {
		c.String(http.StatusOK, "history")
	}))
	r.GET("/records/immovable/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "record")
	})

	tests := []struct {
		path string
		want int
		body string
	}{
		{"/records/immovable/uuid", http.StatusOK, "record"},
		{"/records/immovable/uuid/pdf", http.StatusOK, "pdf"},
		{"/records/immovable/uuid/history", http.StatusOK, "history"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.path, nil)
			r.ServeHTTP(w, req)
			if w.Code != tt.want {
				t.Fatalf("status = %d, want %d", w.Code, tt.want)
			}
			if w.Body.String() != tt.body {
				t.Fatalf("body = %q, want %q", w.Body.String(), tt.body)
			}
		})
	}
}
