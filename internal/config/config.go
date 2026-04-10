package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabasePath   string
	AdminSecretKey string
}

func Load() (*Config, error) {
	// .env file is optional — existing env vars take precedence
	_ = godotenv.Load()

	adminKey := os.Getenv("ADMIN_SECRET_KEY")
	if adminKey == "" {
		return nil, fmt.Errorf("ADMIN_SECRET_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "linkor.db"
	}

	return &Config{
		Port:           port,
		DatabasePath:   dbPath,
		AdminSecretKey: adminKey,
	}, nil
}
