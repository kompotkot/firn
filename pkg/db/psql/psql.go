//go:build psql

package psql

import (
	"context"
	"strings"
	"time"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/kb"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PsqlDB represents a PostgreSQL database connection
type PsqlDB struct {
	pool *pgxpool.Pool
}

// NewPsqlDB creates a new PostgreSQL database connection
func NewPsqlDB(uri string, maxConns int, connMaxLifetime time.Duration) (*PsqlDB, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = int32(maxConns)
	pool.Config().MaxConnLifetime = connMaxLifetime

	return &PsqlDB{pool: pool}, nil
}

// TestConnection tests the database connection with a timeout
func (p *PsqlDB) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.pool.Ping(ctx)
}

// Close closes the database connection pool
func (p *PsqlDB) Close() error {
	if p.pool != nil {
		p.pool.Close()
	}
	return nil
}

// ListJournals lists all journals ordered by updated_at
func (p *PsqlDB) ListJournals(ctx context.Context, orderByDesc bool, limit, offset int) ([]kb.Journal, error) {
	var sb strings.Builder

	sb.WriteString("SELECT id, name, created_at, updated_at FROM journals ORDER BY updated_at")
	if orderByDesc {
		sb.WriteString(" DESC")
	}
	sb.WriteString(" LIMIT $1 OFFSET $2")

	query := sb.String()

	if limit == 0 {
		limit = db.JOURNAL_LIST_DEFAULT_LIMIT
	}

	rows, err := p.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[kb.Journal])
}

// ListEntries lists all entries for a journal
func (p *PsqlDB) ListEntries(ctx context.Context, journalId string, orderByDesc bool, limit, offset int) ([]kb.Entry, error) {
	var sb strings.Builder

	sb.WriteString("SELECT id, journal_id, title, content, created_at, updated_at FROM entries WHERE journal_id = $1 ORDER BY updated_at")
	if orderByDesc {
		sb.WriteString(" DESC")
	}
	sb.WriteString(" LIMIT $2 OFFSET $3")

	query := sb.String()

	if limit == 0 {
		limit = db.ENTRY_LIST_DEFAULT_LIMIT
	}

	rows, err := p.pool.Query(ctx, query, journalId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[kb.Entry])
}

// GetJournalById retrieves a journal by its ID
func (p *PsqlDB) GetJournalById(ctx context.Context, id string) (*kb.Journal, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// CreateJournal creates a new journal with the given name
func (p *PsqlDB) CreateJournal(ctx context.Context, name string) (*kb.Journal, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteJournal deletes a journal by its ID
func (p *PsqlDB) DeleteJournal(ctx context.Context, id string) error {
	// TODO(kompotkot): Implement
	return nil
}

// GetEntryById retrieves an entry by journal ID and entry ID
func (p *PsqlDB) GetEntryById(ctx context.Context, journalId, entryId string) (*kb.Entry, error) {
	query := "SELECT id, journal_id, title, content, created_at, updated_at FROM entries WHERE journal_id = $1 AND id = $2"

	row := p.pool.QueryRow(ctx, query, journalId, entryId)

	var entry kb.Entry
	err := row.Scan(&entry.Id, &entry.JournalId, &entry.Title, &entry.Content, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &entry, nil
}

// CreateEntry creates a new entry in the specified journal
func (p *PsqlDB) CreateEntry(ctx context.Context, journalId, title, content string) (*kb.Entry, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteEntry deletes an entry by journal ID and entry ID
func (p *PsqlDB) DeleteEntry(ctx context.Context, journalId, entryId string) error {
	// TODO(kompotkot): Implement
	return nil
}

// ListEntryTags lists all tags assigned to an entry
func (p *PsqlDB) ListEntryTags(ctx context.Context, journalId, entryId string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// AssignTagsToEntry assigns tags to an entry
func (p *PsqlDB) AssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error {
	// TODO(kompotkot): Implement
	return nil
}

// DeAssignTagsToEntry removes tag assignments from an entry
func (p *PsqlDB) DeAssignTagsToEntry(ctx context.Context, journalId, entryId string, tagIds []string) error {
	// TODO(kompotkot): Implement
	return nil
}

// ListTags lists all tags, optionally filtered by labels
func (p *PsqlDB) ListTags(ctx context.Context, labels []string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// CreateTags creates new tags with the given labels
func (p *PsqlDB) CreateTags(ctx context.Context, labels []string) ([]kb.Tag, error) {
	// TODO(kompotkot): Implement
	return nil, nil
}

// DeleteTags deletes tags by their IDs
func (p *PsqlDB) DeleteTags(ctx context.Context, ids []string) error {
	// TODO(kompotkot): Implement
	return nil
}
