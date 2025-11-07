package db

import (
	"context"

	"github.com/kompotkot/firn/pkg/journal"
)

const JOURNAL_LIST_DEFAULT_LIMIT int = 100
const ENTRY_LIST_DEFAULT_LIMIT int = 100

// Database represents a common interface for database operations
type Database interface {
	// TestConnection tests the database connection with a timeout
	TestConnection(ctx context.Context) error

	// ListJournals lists all journals ordered by updated_at
	ListJournals(ctx context.Context, orderByDesc bool, limit, offset int) ([]journal.Journal, error)

	// ListEntries lists all entries for a journal
	ListEntries(ctx context.Context, journalId string, orderByDesc bool, limit, offset int) ([]journal.Entry, error)

	// Close closes the database connection
	Close() error
}
