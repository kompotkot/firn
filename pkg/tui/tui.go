//go:build tui

package tui

import (
	"context"
	"fmt"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/journal"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	database db.Database

	ready bool

	keys keymap
	help help.Model

	// List of journals
	journalList list.Model

	viewport viewport.Model

	// Window width for dynamic item rendering
	width int

	// Debug
	debugActive bool
	debugStr    string
}

// Initialize TUI model
func initModel(database db.Database) model {
	return model{
		database: database,

		keys: initKeymap(),

		journalList: initJournalList(),

		debugActive: initDebug(),
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
				items[i] = jItem{
					journal:    j,
					widthTitle: m.width - len(j.Name) - magicWidthPaddingNum,
					widthDesc:  m.width - len(fmt.Sprintf("ID: %s", j.Id)) - magicWidthPaddingNum, // Available width for Description (full width minus ID length)
				}
			}
			m.journalList.SetItems(items)
			m.journalList.Paginator.SetTotalPages(len(msg))
		}
		m.viewport.SetContent("No journals found..")

	// Window size changed
	case tea.WindowSizeMsg:
		header := m.headerView()
		footer := m.footerView()

		// Calculate actual heights including any styling (borders, padding, etc.)
		headerHeight := lipgloss.Height(header)
		footerHeight := lipgloss.Height(footer)

		// lipgloss.JoinVertical adds newlines between elements
		// For 3 elements (header, viewport, footer) there will be 2 newlines

		viewportHeight := msg.Height - headerHeight - footerHeight

		// Ensure viewport height is at least 0
		if viewportHeight < 0 {
			viewportHeight = 0
		}

		// Update width for dynamic item rendering
		m.width = msg.Width

		if !m.ready {
			// Initialize viewport only after receiving window dimensions (they arrive asynchronously)
			m.viewport = viewport.New(msg.Width, viewportHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = viewportHeight
		}

		m.journalList.SetSize(msg.Width, viewportHeight)

		// Update all items with new width
		currentItems := m.journalList.Items()
		updatedItems := make([]list.Item, len(currentItems))
		for i, it := range currentItems {
			if ji, ok := it.(jItem); ok {
				ji.widthTitle = m.width - len(ji.journal.Name) - rightPaddingDatetime
				ji.widthDesc = m.width - len(fmt.Sprintf("ID: %s", ji.journal.Id)) - rightPaddingDatetime // Available width for Description (full width minus ID length)
				updatedItems[i] = ji
			} else {
				updatedItems[i] = it
			}
		}
		m.journalList.SetItems(updatedItems)

		// Update viewport content after resize
		m.viewport.SetContent(m.journalList.View())

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
	footer := m.footerView()

	// Build the view: header, viewport, footer
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}

func ShowTui(ctx context.Context, database db.Database) error {
	p := tea.NewProgram(initModel(database), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
