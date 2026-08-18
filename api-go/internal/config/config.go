// Package config loads configuration from environment variables.
package config

import "os"

// Config holds all application configuration.
type Config struct {
	JWTSecret    string
	NodeAPIURL   string
	Port         string
}

// Load reads configuration from environment with sensible defaults.
// DECISION: Using os.LookupEnv to allow overrides while providing defaults.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	nodeURL := os.Getenv("NODE_API_URL")
	if nodeURL == "" {
		nodeURL = "http://localhost:3000"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}
	return &Config{
		JWTSecret:  secret,
		NodeAPIURL: nodeURL,
		Port:       port,
	}
}
