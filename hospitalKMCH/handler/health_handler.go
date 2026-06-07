package handler

import "github.com/gin-gonic/gin"

func HealthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	}
}
