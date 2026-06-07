package conf

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SERVICE_DB_HOST string
	SERVICE_DB_PORT string
	SERVICE_DB_USER string
	SERVICE_DB_PASS string
	SERVICE_DB_NAME string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		SERVICE_DB_HOST: os.Getenv("SERVICE_DB_HOST"),
		SERVICE_DB_USER: os.Getenv("SERVICE_DB_USER"),
		SERVICE_DB_PASS: os.Getenv("SERVICE_DB_PASS"),
		SERVICE_DB_PORT: os.Getenv("SERVICE_DB_PORT"),
		SERVICE_DB_NAME: os.Getenv("SERVICE_DB_NAME"),
	}, nil

}
