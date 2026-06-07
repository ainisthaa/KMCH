package conf

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_NAME    string
	APP_PORT    string

	SERVICE_DB_HOST string
	SERVICE_DB_PORT string
	SERVICE_DB_USER string
	SERVICE_DB_PASS string
	SERVICE_DB_NAME string

	EXCEL_FILE_PATH          string // registrations pre-check Excel
	MENTAL_HEALTH_EXCEL_PATH string // mental-health screening completed Excel
	MENTAL_HEALTH_ISSUE_PATH string // mental-health issue/risk Excel (separate file)
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()
	return &Config{
		APP_NAME:    getEnv("APP_NAME", "lineoa-miniapp"),
		APP_PORT:    getEnv("APP_PORT", "8080"),

		SERVICE_DB_HOST: getEnv("SERVICE_DB_HOST", "localhost"),
		SERVICE_DB_PORT: getEnv("SERVICE_DB_PORT", "3306"),
		SERVICE_DB_USER: getEnv("SERVICE_DB_USER", "appuser"),
		SERVICE_DB_PASS: getEnv("SERVICE_DB_PASS", "password"),
		SERVICE_DB_NAME: getEnv("SERVICE_DB_NAME", "lineoa_miniapp"),

		EXCEL_FILE_PATH:          getEnv("EXCEL_FILE_PATH", "./data/registrations.xlsx"),
		MENTAL_HEALTH_EXCEL_PATH: getEnv("MENTAL_HEALTH_EXCEL_PATH", "./data/mental_health.xlsx"),
		MENTAL_HEALTH_ISSUE_PATH: getEnv("MENTAL_HEALTH_ISSUE_PATH", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
