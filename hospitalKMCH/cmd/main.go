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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"lineoa-miniapp/conf"
	"lineoa-miniapp/domain"
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
	} else {
		if err := autoMigrate(db); err != nil {
			applog.UnhandledError("auto_migrate", err)
		}
		if err := seedDefaults(db); err != nil {
			applog.UnhandledError("seed_defaults", err)
		}
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
	displaySvc := service.NewDisplayService(roomRepo, queueRepo, patientRepo)
	staffSvc := service.NewStaffService(roomRepo, queueRepo, patientRepo)

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

func autoMigrate(db *gorm.DB) error {
	applog.Log.Info().Str("action", "auto_migrate").Msg("running gorm auto-migrate")
	return db.AutoMigrate(
		&domain.EventInfo{},
		&domain.PatientInfo{},
		&domain.DoctorRoom{},
		&domain.PatientCheck{},
		&domain.PatientQueue{},
	)
}

// seedDefaults makes sure the minimum data the API needs is present: the
// default event row (id=1) and the five doctor rooms. Idempotent.
func seedDefaults(db *gorm.DB) error {
	now := time.Now()
	event := domain.EventInfo{
		EventID:       1,
		EventName:     "Default Event",
		EventDateFrom: now,
		EventDateTo:   now.Add(24 * time.Hour),
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&event).Error; err != nil {
		return err
	}

	rooms := []domain.DoctorRoom{
		{RoomID: "room-001", RoomName: "Room 1"},
		{RoomID: "room-002", RoomName: "Room 2"},
		{RoomID: "room-003", RoomName: "Room 3"},
		{RoomID: "room-004", RoomName: "Room 4"},
		{RoomID: "room-005", RoomName: "Room 5"},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rooms).Error; err != nil {
		return err
	}
	return nil
}
