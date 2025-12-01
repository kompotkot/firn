package db

import (
	"context"

	"github.com/kompotkot/firn/pkg/kb"
)

const JOURNAL_LIST_DEFAULT_LIMIT int = 100
const ENTRY_LIST_DEFAULT_LIMIT int = 100

// Database represents a common interface for database operations
type Database interface {
	// TestConnection tests the database connection with a timeout
	TestConnection(ctx context.Context) error

	// ListJournals lists all journals ordered by updated_at
	ListJournals(ctx context.Context, orderByDesc bool, limit, offset int) ([]kb.Journal, error)

	// ListEntries lists all entries for a journal
	ListEntries(ctx context.Context, journalId string, orderByDesc bool, limit, offset int) ([]kb.Entry, error)

	// GetJournalById retrieves a journal by its ID
	GetJournalById(ctx context.Context, id string) (*kb.Journal, error)

	// CreateJournal creates a new journal with the given name
	CreateJournal(ctx context.Context, name string) (*kb.Journal, error)

	// DeleteJournal deletes a journal by its ID
	DeleteJournal(ctx context.Context, id string) error

	// GetEntryById retrieves an entry by journal ID and entry ID
	GetEntryById(ctx context.Context, journalId, entryId string) (*kb.Entry, error)

	// CreateEntry creates a new entry in the specified journal
	CreateEntry(ctx context.Context, journalId, title, content string) (*kb.Entry, error)

	// DeleteEntry deletes an entry by journal ID and entry ID
	DeleteEntry(ctx context.Context, journalId, entryId string) error

	// ListEntryTags lists all tags assigned to an entry
	ListEntryTags(ctx context.Context, journalId, entryId string) ([]kb.Tag, error)

	// AssignTagsToEntry assigns tags to an entry
	AssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error

	// DeAssignTagsToEntry removes tag assignments from an entry
	DeAssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error

	// ListTags lists all tags, optionally filtered by labels
	ListTags(ctx context.Context, labels []string) ([]kb.Tag, error)

	// CreateTags creates new tags with the given labels
	CreateTags(ctx context.Context, labels []string) ([]kb.Tag, error)

	// DeleteTags deletes tags by their IDs
	DeleteTags(ctx context.Context, ids []string) error

	// Close closes the database connection
	Close() error
}
