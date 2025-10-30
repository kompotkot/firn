package types

import "time"

// Logger configuration
type LoggerConfig struct {
	Level  string
	Format string
}

// Database configuration
type DatabaseConfig struct {
	Type            string
	URI             string
	MaxConns        int
	ConnMaxLifetime time.Duration
}

// Main configuration
type Config struct {
	Logger   LoggerConfig
	Database DatabaseConfig
}
