//go:build psql

package psql

import (
	"context"
	"time"

	"github.com/kompotkot/firn/pkg/journal"

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


// ListJournals lists all journals
func (p *PsqlDB) ListJournals(ctx context.Context) ([]journal.Journal, error) {
	query := `
		SELECT id, name, created_at, updated_at FROM journal
	`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the rows into the journals slice
	var journals []journal.Journal
	for rows.Next() {
		var journal journal.Journal
		err := rows.Scan(&journal.Id, &journal.Name, &journal.CreatedAt, &journal.UpdatedAt)
		if err != nil {
			return nil, err
		}
		journals = append(journals, journal)
	}
	
	// Check for errors from the rows iteration cycle
	if err := rows.Err(); err != nil {
		return nil, err
	}


	return journals, nil
}

// ListEntries lists all entries for a journal
func (p *PsqlDB) ListEntries(ctx context.Context, journalId string) ([]journal.Entry, error) {
	query := `
		SELECT id, journal_id, title, content, created_at, updated_at FROM entry
	`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the rows into the entries slice
	var entries []journal.Entry
	for rows.Next() {
		var entry journal.Entry
		err := rows.Scan(&entry.Id, &entry.JournalId, &entry.Title, &entry.Content, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	
	// Check for errors from the rows iteration cycle
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
