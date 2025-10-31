//go:build tui

package tui

import (
	"context"

	"github.com/kompotkot/firn/pkg/db"

	tea "github.com/charmbracelet/bubbletea"
)

// List journals from the database and return as tea data
func listJournals(database db.Database) tea.Cmd {
	return func() tea.Msg {
		journals, err := database.ListJournals(context.Background())
		if err != nil {
			return err
		}
		return journals
	}
}

