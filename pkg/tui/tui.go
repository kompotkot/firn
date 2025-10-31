//go:build tui

package tui

import (
	"context"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/journal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

type model struct {
	database db.Database

	journals []journal.Journal
}

// Initialize TUI model
func initModel(database db.Database) model {
	return model{
		database: database,

		journals: []journal.Journal{},
	}
}

// Execute commands concurrently with no ordering guarantees during initialization
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listJournals(m.database),
	)
}

// Processes events like window resize, errors, loaded data, and key presses
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

// Assembles the UI string for each frame
func (m model) View() string {
	return ""
}

func ShowTui(ctx context.Context, database db.Database) error {
	p := tea.NewProgram(initModel(database), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
