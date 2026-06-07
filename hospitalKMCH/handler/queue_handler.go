package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/dto"
	"lineoa-miniapp/service"
)

func ScanDoctorQueue(svc service.QueueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		var req dto.ScanDoctorQueueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		if err := svc.ScanDoctorQueue(c.Request.Context(), lineID, req.EventID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully joined the doctor consultation queue."})
	}
}

func GetQueueStatus(svc service.QueueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		eventID, _ := strconv.ParseUint(c.Query("event_id"), 10, 64)
		if eventID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "event_id query param required"})
			return
		}
		resp, err := svc.GetQueueStatus(c.Request.Context(), lineID, uint(eventID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CompleteDoctorConsultation(svc service.QueueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lineID := c.Param("line_id")
		var req dto.CompleteConsultationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		if err := svc.CompleteDoctorConsultation(c.Request.Context(), lineID, req.EventID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":      "Doctor consultation completed.",
			"next_station": "xray",
		})
	}
}
