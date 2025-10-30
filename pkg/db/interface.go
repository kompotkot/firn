package db

import (
	"context"
)

// Database represents a common interface for database operations
type Database interface {
	// TestConnection tests the database connection with a timeout
	TestConnection(ctx context.Context) error

	// Close closes the database connection
	Close() error
}
