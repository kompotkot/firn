package db

import (
	"context"

	"github.com/kompotkot/firn/pkg/journal"
)

// Database represents a common interface for database operations
type Database interface {
	// TestConnection tests the database connection with a timeout
	TestConnection(ctx context.Context) error

	// ListJournals lists all journals
	ListJournals(ctx context.Context) ([]journal.Journal, error)

	// ListEntries lists all entries for a journal
	ListEntries(ctx context.Context, journalId string) ([]journal.Entry, error)

	// Close closes the database connection
	Close() error
}
