package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration values
type Config struct {
	PortNumber       string
	CertFile         string
	KeyFile          string
	Database         DatabaseConfig
	MessagesDatabase DatabaseConfig
}

// DatabaseConfig holds the configuration for the database
type DatabaseConfig struct {
	Path       string
	InitScript string
}

// AppConfig holds the global application configuration
var AppConfig Config

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	AppConfig = Config{
		PortNumber:       mustGetEnv("PORT_NUMBER"),
		CertFile:         mustGetEnv("CERT_FILE"),
		KeyFile:          mustGetEnv("KEY_FILE"),
		Database:         loadDatabaseConfig("DB_PATH", "DB_INIT_SCRIPT"),
		MessagesDatabase: loadDatabaseConfig("MESSAGES_DB_PATH", ""),
	}
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig(pathKey, initScriptKey string) DatabaseConfig {
	return DatabaseConfig{
		Path:       mustGetEnv(pathKey),
		InitScript: getEnv(initScriptKey, ""),
	}
}

// mustGetEnv reads an environment variable or logs a fatal error if not found
func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}