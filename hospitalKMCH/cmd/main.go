package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"lineoa-miniapp/conf"
	"lineoa-miniapp/internal/dbconn"
	applog "lineoa-miniapp/pkg/logger"
	"lineoa-miniapp/pkg/mentalhealthcache"
	excelrepo "lineoa-miniapp/repository/excel"
	mysqlrepo "lineoa-miniapp/repository/mysql"
	"lineoa-miniapp/routes"
	"lineoa-miniapp/service"
)

func main() {
	applog.Init("/logs/app.log")
	applog.Log.Info().Str("action", "startup").Msg("starting lineoa-miniapp")

	config, err := conf.NewConfig()
	if err != nil {
		applog.UnhandledError("load_config", err)
		os.Exit(1)
	}
	applog.Log.Info().
		Str("action", "config_loaded").
		Str("app_name", config.APP_NAME).
		Str("port", config.APP_PORT).
		Msg("config loaded")

	db, dbErr := dbconn.TryDBConnect(
		config.SERVICE_DB_USER, config.SERVICE_DB_PASS,
		config.SERVICE_DB_HOST, config.SERVICE_DB_PORT, config.SERVICE_DB_NAME,
	)
	if dbErr != nil {
		applog.Log.Warn().Err(dbErr).Str("action", "db_connect").Msg("database unavailable at startup")
	}

	issuePath := config.MENTAL_HEALTH_ISSUE_PATH
	if issuePath == "" {
		issuePath = config.MENTAL_HEALTH_EXCEL_PATH
	}
	mhCache := mentalhealthcache.NewCache(config.MENTAL_HEALTH_EXCEL_PATH)
	_ = issuePath

	// ── Repositories ──────────────────────────────────────────────────────────
	registrationRepo := excelrepo.NewExcelRegistrationRepository(config.EXCEL_FILE_PATH)
	patientRepo := mysqlrepo.NewMySQLPatientRepository(db)
	roomRepo := mysqlrepo.NewMySQLRoomRepository(db)
	queueRepo := mysqlrepo.NewMySQLQueueRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	registerSvc := service.NewRegisterService(patientRepo, registrationRepo, mhCache)
	routeSvc := service.NewRouteService(db, patientRepo, mhCache)
	queueSvc := service.NewQueueService(db, patientRepo, queueRepo, roomRepo)
	displaySvc := service.NewDisplayService(roomRepo, queueRepo)
	staffSvc := service.NewStaffService(roomRepo, queueRepo)

	// ── Startup room fill ─────────────────────────────────────────────────────
	go func() {
		if err := queueSvc.AutoFillRooms(context.Background()); err != nil {
			applog.UnhandledError("startup_auto_fill", err)
		}
	}()

	// ── HTTP server ───────────────────────────────────────────────────────────
	server := gin.Default()
	routes.RegisterRoutes(server, registerSvc, routeSvc, queueSvc, displaySvc, staffSvc)

	addr := fmt.Sprintf(":%s", config.APP_PORT)
	srv := &http.Server{Addr: addr, Handler: server}

	go func() {
		applog.Log.Info().Str("action", "server_start").Str("addr", addr).Msg("server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			applog.UnhandledError("server_listen", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	applog.Log.Info().Str("action", "shutdown").Msg("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		applog.UnhandledError("server_shutdown", err)
	}
	applog.Log.Info().Str("action", "stopped").Msg("server stopped")
}
