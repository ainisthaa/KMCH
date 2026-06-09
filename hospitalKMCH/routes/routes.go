package routes

import (
	"github.com/gin-gonic/gin"
	"lineoa-miniapp/handler"
	"lineoa-miniapp/middleware"
	"lineoa-miniapp/service"
)

func RegisterRoutes(
	server *gin.Engine,
	registerSvc service.RegisterService,
	routeSvc    service.RouteService,
	queueSvc    service.QueueService,
	displaySvc  service.DisplayService,
	staffSvc    service.StaffService,
) {
	server.Use(middleware.RequestID())
	server.Use(middleware.Logger())
	server.Use(middleware.Recovery())
	server.Use(middleware.CORS())

	server.GET("/health", handler.HealthCheck())

	api := server.Group("/api/v1")
	{
		api.GET("/health", handler.HealthCheck())

		// ── Registration ────────────────────────────────────────────────────
		api.POST("/register", handler.Register(registerSvc))

		// ── Patient journey ─────────────────────────────────────────────────
		p := api.Group("/patients/:line_id")
		{
			p.POST("/scan-after-payment",           handler.ScanAfterPayment(routeSvc))
			p.GET("/route",                         handler.GetRoute(routeSvc))
			p.GET("/check",                         handler.GetCheckStatus(routeSvc))
			p.POST("/scan-doctor-queue",            handler.ScanDoctorQueue(queueSvc))
			p.GET("/queue-status",                  handler.GetQueueStatus(queueSvc))
			p.POST("/complete-doctor-consultation", handler.CompleteDoctorConsultation(queueSvc))
		}

		// ── Public display board ────────────────────────────────────────────
		api.GET("/display", handler.GetDisplay(displaySvc))

		// ── Staff ───────────────────────────────────────────────────────────
		staff := api.Group("/staff")
		{
			staff.GET("/dashboard",             handler.StaffDashboard(staffSvc))
			staff.GET("/queue",                 handler.StaffQueue(staffSvc))
			staff.POST("/queue/:queue_id/skip", handler.SkipQueue(queueSvc))
			staff.POST("/fill",                 handler.StaffFill(queueSvc))
		}
	}
}
