//go:build tui

package tui

import (
	"context"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/journal"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
)

// Keymap for help panel
type keymap = struct {
	quit key.Binding
}

type model struct {
	database db.Database

	ready bool

	keys keymap
	help help.Model

	journals []journal.Journal

	viewport viewport.Model
}

// Initialize TUI model
func initModel(database db.Database) model {
	return model{
		database: database,

		keys: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},

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
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		// Handle exit from TUI
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		header := m.headerView()
		help := m.help.ShortHelpView([]key.Binding{m.keys.quit})

		// Calculate actual heights including any styling
		headerHeight := lipgloss.Height(header)
		helpHeight := lipgloss.Height(help)

		// lipgloss.JoinVertical adds newlines between elements
		// For 3 elements (header, viewport, help) there will be 2 newlines
		verticalMargin := 0 // TODO(kompotkot): Check if set to 2 when content is added
		viewportHeight := msg.Height - headerHeight - helpHeight - verticalMargin

		// Ensure viewport height is at least 0
		if viewportHeight < 0 {
			viewportHeight = 0
		}

		if !m.ready {
			// Initialize viewport only after receiving window dimensions (they arrive asynchronously)
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
		}
	}

	// Handle keyboard events in the viewport
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// Assembles the UI string for each frame
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	header := m.headerView()
	help := m.help.ShortHelpView([]key.Binding{m.keys.quit})

	// Build the view: header, viewport, help
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		help,
	)
}

func ShowTui(ctx context.Context, database db.Database) error {
	p := tea.NewProgram(initModel(database), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
