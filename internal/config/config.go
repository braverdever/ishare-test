package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string
	Issuer     string
	Audience   string
	Expiration time.Duration
}

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Environment string
	Port        string
}

// Load loads configuration from environment variables
func Load() *Config {
	expiration, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "ishare_tasks"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			Issuer:     getEnv("JWT_ISSUER", "ishare-task-api"),
			Audience:   getEnv("JWT_AUDIENCE", "ishare-clients"),
			Expiration: time.Duration(expiration) * time.Hour,
		},
		OAuth: OAuthConfig{
			ClientID:     getEnv("OAUTH_CLIENT_ID", "test-client"),
			ClientSecret: getEnv("OAUTH_CLIENT_SECRET", "test-secret"),
			RedirectURI:  getEnv("OAUTH_REDIRECT_URI", "http://localhost:8080/oauth/callback"),
		},
		Server: ServerConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnv("SERVER_PORT", "8080"),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 