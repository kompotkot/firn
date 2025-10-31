package db

import "errors"

var (
	ErrJournalNotFound = errors.New("journal not found")
	ErrEntryNotFound   = errors.New("entry not found")
)
