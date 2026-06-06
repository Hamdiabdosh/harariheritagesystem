package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		log.Printf("[%s] %d | %13v | %15s | %s",
			method, status, latency, clientIP, path,
		)
	}
}
