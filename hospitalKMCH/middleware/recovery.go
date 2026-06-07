package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.Printf("[RECOVERY] panic recovered: %v", err)
		c.JSON(500, gin.H{"message": "internal server error"})
	})
}
