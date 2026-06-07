package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/service"
)

func StaffDashboard(svc service.StaffService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := svc.GetDashboard(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func StaffQueue(svc service.StaffService) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := svc.GetWaitingQueue(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func SkipQueue(svc service.QueueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		queueID, err := strconv.ParseUint(c.Param("queue_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "invalid queue_id"})
			return
		}
		if err := svc.SkipQueue(c.Request.Context(), uint(queueID)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Queue entry skipped."})
	}
}

func StaffFill(svc service.QueueService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := svc.AutoFillRooms(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room fill triggered."})
	}
}
