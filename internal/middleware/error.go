package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if c.Writer.Written() {
			log.Printf("error after response written: %v", err)
			return
		}

		respondError(c, http.StatusInternalServerError, "Internal server error")
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Printf("panic recovered: %v", recovered)
		respondError(c, http.StatusInternalServerError, "Internal server error")
	})
}

func respondError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"success": false,
		"error":   message,
		"code":    code,
	})
}

func RespondSuccess(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"message": message,
	})
}

func RespondError(c *gin.Context, code int, message string) {
	respondError(c, code, message)
}
