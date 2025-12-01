package kb

import "time"

// Journal represents a journal in the system
type Journal struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Entry represents an entry in a journal
type Entry struct {
	Id        string    `json:"id"`
	JournalId string    `json:"journal_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Tag represents a label assigned to journal entry
type Tag struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

// Assignments of tags to entries
type TagAssignment struct {
	TagId   string `json:"tag_id"`
	EntryId string `json:"entry_id"`
}
