package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/service"
)

func GetDisplay(svc service.DisplayService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := svc.GetDisplay(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
