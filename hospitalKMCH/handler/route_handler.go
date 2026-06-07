package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/dto"
	"lineoa-miniapp/service"
)

func ScanAfterPayment(svc service.RouteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		var req dto.ScanAfterPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		resp, err := svc.ScanAfterPayment(c.Request.Context(), lineID, req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func GetRoute(svc service.RouteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		eventID, _ := strconv.ParseUint(c.Query("event_id"), 10, 64)
		if eventID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "event_id query param required"})
			return
		}
		resp, err := svc.GetRoute(c.Request.Context(), lineID, uint(eventID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CompletePsychologist(svc service.RouteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		var req dto.CompleteStationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		resp, err := svc.CompletePsychologist(c.Request.Context(), lineID, req.EventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CompleteRightsTransfer(svc service.RouteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		var req dto.CompleteStationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		resp, err := svc.CompleteRightsTransfer(c.Request.Context(), lineID, req.EventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
