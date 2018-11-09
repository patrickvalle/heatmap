package config

import (
	"os"
	"time"
)

// Config holds the app configuration.
type Config struct {
	APIHost         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// New creates and returns a new Config.
func New() Config {

	// Set the default values.
	cfg := Config{
		APIHost:         ":8080",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}

	// Attempt to populate values from environment variables.
	if apiHost := os.Getenv("API_HOST"); apiHost != "" {
		cfg.APIHost = apiHost
	}

	return cfg
}
