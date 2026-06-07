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
		method := c.Request.Method

		log.Printf("[REQUEST] %s %s from %s", method, path, c.ClientIP())

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		log.Printf("[RESPONSE] %s %s status=%d duration=%s", method, path, status, duration)
	}
}
