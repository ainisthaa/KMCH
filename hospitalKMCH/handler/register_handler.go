package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/dto"
	"lineoa-miniapp/service"
)

func Register(svc service.RegisterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		resp, err := svc.Register(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
