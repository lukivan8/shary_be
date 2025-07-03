package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Port        int
	DatabaseURL string
	Environment string
	LogLevel    string
}

// Load loads configuration from environment variables and .env file
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we'll use system environment variables
		log.Println("No .env file found, using system environment variables")
	}

	port := 4000 // default port
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Default local PostgreSQL connection string
		databaseURL = "postgres://postgres:password@localhost:5432/shary_be?sslmode=disable"
	} else {
		// Ensure sslmode=disable is present in the DATABASE_URL
		if !contains(databaseURL, "sslmode=") {
			if contains(databaseURL, "?") {
				databaseURL += "&sslmode=disable"
			} else {
				databaseURL += "?sslmode=disable"
			}
		}
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
		Environment: environment,
		LogLevel:    logLevel,
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
