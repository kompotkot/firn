//go:build sqlite

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/kb"

	_ "github.com/mattn/go-sqlite3"
)

// validSyncModes lists the allowed values for the synchronous pragma
var validSyncModes = map[string]bool{
	"OFF":    true,
	"NORMAL": true,
	"FULL":   true,
	"EXTRA":  true,
}

// SqliteDB represents a SQLite database connection
type SqliteDB struct {
	db *sql.DB
}

// NewSqliteDB creates a new SQLite database connection with specified options
func NewSqliteDB(uri string, enableWal bool, syncPragma string) (*SqliteDB, error) {
	params := url.Values{}

	if enableWal {
		params.Add("_journal_mode", "WAL")
	}

	if syncPragma != "" {
		ucSyncPragma := strings.ToUpper(syncPragma)
		if !validSyncModes[ucSyncPragma] {
			return nil, fmt.Errorf("invalid sync pragma value: %s. Must be one of OFF, NORMAL, FULL, EXTRA", syncPragma)
		}
		params.Add("_synchronous", ucSyncPragma)
	}

	constructedUri := uri
	if len(params) > 0 {
		if strings.Contains(uri, "?") {
			constructedUri += "&" + params.Encode()
		} else {
			constructedUri += "?" + params.Encode()
		}
	}

	db, err := sql.Open("sqlite3", constructedUri)
	if err != nil {
		return nil, fmt.Errorf("failed to open database with DSN '%s': %w", constructedUri, err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// Enable foreign key support for this connection.
	// This is crucial for ON DELETE CASCADE and other FK actions to work.
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close() // Close DB if we can't set the pragma
		return nil, fmt.Errorf("failed to enable foreign key support for DSN '%s': %w", constructedUri, err)
	}

	return &SqliteDB{db: db}, nil
}

// TestConnection tests the database connection with a timeout
func (s *SqliteDB) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

// Close closes the database connection
func (s *SqliteDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// ListJournals lists all journals ordered by updated_at
func (s *SqliteDB) ListJournals(ctx context.Context, orderByDesc bool, limit, offset int) ([]kb.Journal, error) {
	var sb strings.Builder

	sb.WriteString("SELECT id, name, created_at, updated_at FROM journals ORDER BY updated_at")
	if orderByDesc {
		sb.WriteString(" DESC")
	}
	sb.WriteString(" LIMIT ? OFFSET ?")

	query := sb.String()

	if limit == 0 {
		limit = db.JOURNAL_LIST_DEFAULT_LIMIT
	}

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var journals []kb.Journal
	for rows.Next() {
		var j kb.Journal
		if err := rows.Scan(&j.Id, &j.Name, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		journals = append(journals, j)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}

// ListEntries lists all entries for a journal
func (s *SqliteDB) ListEntries(ctx context.Context, journalId string, orderByDesc bool, limit, offset int) ([]kb.Entry, error) {
	var sb strings.Builder

	sb.WriteString("SELECT id, journal_id, title, content, created_at, updated_at FROM entries WHERE journal_id = ? ORDER BY updated_at")
	if orderByDesc {
		sb.WriteString(" DESC")
	}
	sb.WriteString(" LIMIT ? OFFSET ?")

	query := sb.String()

	if limit == 0 {
		limit = db.ENTRY_LIST_DEFAULT_LIMIT
	}

	rows, err := s.db.QueryContext(ctx, query, journalId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []kb.Entry
	for rows.Next() {
		var e kb.Entry
		if err := rows.Scan(&e.Id, &e.JournalId, &e.Title, &e.Content, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetJournalById retrieves a journal by its ID
func (s *SqliteDB) GetJournalById(ctx context.Context, id string) (*kb.Journal, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// CreateJournal creates a new journal with the given name
func (s *SqliteDB) CreateJournal(ctx context.Context, name string) (*kb.Journal, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteJournal deletes a journal by its ID
func (s *SqliteDB) DeleteJournal(ctx context.Context, id string) error {
	// TODO(kompotkot): Implement
	return nil
}

// GetEntryById retrieves an entry by journal ID and entry ID
func (s *SqliteDB) GetEntryById(ctx context.Context, journalId, entryId string) (*kb.Entry, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// CreateEntry creates a new entry in the specified journal
func (s *SqliteDB) CreateEntry(ctx context.Context, journalId, title, content string) (*kb.Entry, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteEntry deletes an entry by journal ID and entry ID
func (s *SqliteDB) DeleteEntry(ctx context.Context, journalId, entryId string) error {
	// TODO(kompotkot): Implement
	return nil
}

// ListEntryTags lists all tags assigned to an entry
func (s *SqliteDB) ListEntryTags(ctx context.Context, journalId, entryId string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// AssignTagsToEntry assigns tags to an entry
func (s *SqliteDB) AssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error {
	// TODO(kompotkot): Implement
	return nil
}

// DeAssignTagsToEntry removes tag assignments from an entry
func (s *SqliteDB) DeAssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error {
	// TODO(kompotkot): Implement
	return nil
}

// ListTags lists all tags, optionally filtered by labels
func (s *SqliteDB) ListTags(ctx context.Context, labels []string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// CreateTags creates new tags with the given labels
func (s *SqliteDB) CreateTags(ctx context.Context, labels []string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteTags deletes tags by their IDs
func (s *SqliteDB) DeleteTags(ctx context.Context, ids []string) error {
	// TODO(kompotkot): Implement
	return nil
}
