//go:build tui

package tui

import (
	"context"
	"fmt"
	"time"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/journal"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Keymap for help panel
type keymap = struct {
	quit key.Binding
}

// item represents a journal in the list
type item struct {
	journal journal.Journal
}

func (i item) Title() string {
	name := i.journal.Name
	updatedAt := i.journal.UpdatedAt.Format(time.RFC3339)
	return fmt.Sprintf("%s %s", name, updatedAt)
}

func (i item) Description() string {
	id := i.journal.Id
	createdAt := i.journal.CreatedAt.Format(time.RFC3339)
	return fmt.Sprintf("ID: %s, Created At: %s", id, createdAt)
}

func (i item) FilterValue() string { return i.journal.Name }

type model struct {
	database db.Database

	ready bool

	keys keymap
	help help.Model

	// List of journals
	journalList list.Model

	viewport viewport.Model
}

// Initialize TUI model
func initModel(database db.Database) model {
	jld := list.NewDefaultDelegate()
	jld.Styles.SelectedTitle = listSelectedStyle

	// Initialize journals list (dimensions will be set when window size is received)
	jl := list.New([]list.Item{}, jld, 0, 0)
	jl.Title = "Journals"
	jl.SetFilteringEnabled(false)
	jl.SetShowPagination(true)
	jl.Paginator.Type = paginator.Arabic

	return model{
		database: database,

		keys: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
		},

		journalList: jl,
	}
}

// Execute commands concurrently with no ordering guarantees during initialization
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listJournals(m.database, false, 0, 0),
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

	// Journals loaded from database
	case []journal.Journal:
		if len(msg) > 0 {
			items := make([]list.Item, len(msg))
			for i, j := range msg {
				items[i] = item{journal: j}
			}
			m.journalList.SetItems(items)
			m.journalList.Paginator.SetTotalPages(len(msg))
		}
		m.viewport.SetContent("No journals found..")

	// Window size changed
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
			m.journalList.SetSize(msg.Width, viewportHeight)

			m.viewport.SetContent(m.journalList.View())

			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
			m.journalList.SetSize(msg.Width, viewportHeight)

			// Update viewport content after resize
			m.viewport.SetContent(m.journalList.View())
		}

	case error:
		m.viewport.SetContent(fmt.Sprintf("Error occured: %v", msg))
	}

	// Handle viewport events (scrolling, etc.)
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	m.journalList.Paginator, cmd = m.journalList.Paginator.Update(msg)
	cmds = append(cmds, cmd)

	// Handle keyboard events in the list and update viewport content
	if m.ready {
		var listCmd tea.Cmd
		m.journalList, listCmd = m.journalList.Update(msg)
		cmds = append(cmds, listCmd)

		// Update viewport content with the list after all updates
		m.viewport.SetContent(m.journalList.View())
	}

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
