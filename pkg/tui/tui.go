//go:build tui

package tui

import (
	"context"
	"fmt"

	"github.com/kompotkot/firn/pkg/db"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusState string

const (
	focusJournals focusState = "journals"
	focusEntries  focusState = "entries"
	focusEntry    focusState = "entry"
)

type model struct {
	ctx      context.Context
	database db.Database

	ready         bool
	width         int
	contentHeight int

	keys keymap

	selectedJournalId       string
	selectedJournalName     string
	lastJournalIndex        int
	restoreJournalSelection bool

	// Focus state: journals list, entries list, or entry textarea
	focusState focusState

	journals journalsPane
	entries  entriesPane
	viewer   entryViewer

	// Selected entry
	selectedEntryId string
	viewerFull      bool

	// Debug
	debugActive bool
	debugStr    string
}

// Initialize TUI model
func initModel(ctx context.Context, database db.Database) model {
	if ctx == nil {
		ctx = context.Background()
	}

	return model{
		ctx:      ctx,
		database: database,

		keys: initKeymap(),

		journals: newJournalsPane(),
		entries:  newEntriesPane(),
		viewer:   newEntryViewer(),

		focusState:       focusJournals, // Start with journal list focused
		debugActive:      initDebug(),
		lastJournalIndex: -1,
	}
}

// Execute commands concurrently with no ordering guarantees during initialization
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listJournals(m.ctx, m.database, false, 0, 0),
	)
}

// Processes events like window resize, errors, loaded data, and key presses
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	skipListUpdate := false

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
			if m.focusState == focusJournals {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keys.esc):
			switch m.focusState {
			case focusEntries:
				if m.selectedJournalId != "" {
					skipListUpdate = true
					m.selectedJournalId = ""
					m.selectedJournalName = ""
					m.selectedEntryId = ""
					m.setFocusState(focusJournals)
					m.restoreJournalSelection = true

					m.entries.SetTitle("Entries")
					m.entries.Clear()
					m.viewer.SetContent("")

					m.resizeComponents()
				}
			case focusEntry:
				skipListUpdate = true
				m.setFocusState(focusEntries)
				m.resizeComponents()

				// TODO(kompotkot) Enter will trigger save logic later
			}
		case key.Matches(msg, m.keys.enter):
			switch m.focusState {
			case focusJournals:
				if selectedJournal, ok := m.journals.SelectedJournal(); ok {
					m.selectedJournalId = selectedJournal.Id
					m.selectedJournalName = selectedJournal.Name
					m.lastJournalIndex = m.journals.GlobalIndex()
					m.selectedEntryId = ""
					m.setFocusState(focusEntries)

					m.entries.Clear()
					m.entries.SetTitle(fmt.Sprintf("%s entries", m.selectedJournalName))
					m.viewer.SetContent("")
					m.resizeComponents()

					return m, tea.Batch(
						listEntries(m.ctx, m.database, m.selectedJournalId, false, 0, 0),
					)
				}
			case focusEntries:
				if selectedEntry, ok := m.entries.SelectedEntry(); ok {
					m.selectedEntryId = selectedEntry.Id
					m.setFocusState(focusEntry)
					m.resizeComponents()

					return m, getEntryById(m.ctx, m.database, m.selectedJournalId, m.selectedEntryId)
				}
			}
		}

	// Journals loaded from database
	case journalsLoadedMsg:
		if len(msg.journals) > 0 {
			items := make([]list.Item, len(msg.journals))
			for i, j := range msg.journals {
				items[i] = jItem{
					journal:    j,
					widthTitle: m.width - len(j.Name) - magicWidthPaddingNum,
					widthDesc:  m.width - len(fmt.Sprintf("ID: %s", j.Id)) - magicWidthPaddingNum, // Available width for Description (full width minus ID length)
				}
			}
			m.journals.SetItems(items)
			m.journals.SetTotalPages(len(msg.journals))
		}
		m.resizeComponents()

	case entriesLoadedMsg:
		if msg.journalId != m.selectedJournalId {
			break
		}

		// Always update the entry list, even if empty (to clear old entries)
		items := make([]list.Item, len(msg.entries))
		for i, e := range msg.entries {
			idLabel := fmt.Sprintf("ID: %s", e.Id)
			items[i] = eItem{
				entry:      e,
				widthTitle: m.width - len(e.Title) - magicWidthPaddingNum,
				widthDesc:  m.width - len(idLabel) - magicWidthPaddingNum,
			}
		}
		m.entries.SetItems(items)
		m.entries.SetTotalPages(len(msg.entries))

		// Auto-select first entry and load its content if entries exist
		if len(msg.entries) > 0 && m.entries.SelectedItem() != nil {
			selectedEntry := m.entries.SelectedItem().(eItem).entry
			m.selectedEntryId = selectedEntry.Id
			cmds = append(cmds, getEntryById(m.ctx, m.database, m.selectedJournalId, m.selectedEntryId))
		} else {
			// No entries, clear selection
			m.selectedEntryId = ""
			m.viewer.SetContent("")
			if m.focusState == focusEntry {
				m.setFocusState(focusEntries)
			}
		}
		m.resizeComponents()

	case entryLoadedMsg:
		if msg.journalId != m.selectedJournalId || msg.entryId != m.selectedEntryId {
			break
		}

		// Entry loaded by ID - update textarea with content
		if msg.entry != nil {
			// Set content and ensure textarea is updated
			m.viewer.SetContent(msg.entry.Content)
		} else {
			// Entry not found - clear textarea
			m.viewer.SetContent("Entry not found")
		}

		m.ensureTextareaFocus(m.focusState == focusEntry)

		m.resizeComponents()

	// Window size changed
	case tea.WindowSizeMsg:
		// Update width for dynamic item rendering first (needed for header/footer calculation)
		m.width = msg.Width

		header := m.headerView()
		footer := m.footerView()

		// Calculate actual heights including any styling (borders, padding, etc.)
		headerHeight := lipgloss.Height(header)
		footerHeight := lipgloss.Height(footer)

		// Calculate available height for lists
		availableHeight := msg.Height - headerHeight - footerHeight
		if availableHeight < 0 {
			availableHeight = 0
		}

		m.contentHeight = availableHeight
		m.ready = true

		m.resizeComponents()

	case errMsg:
		// Handle errors - could display in a status bar or log
		m.debugStr = fmt.Sprintf("%s: %v", msg.operation, msg.err)
	}

	// Handle keyboard events for lists and textarea
	if m.ready && !skipListUpdate {
		var listCmd tea.Cmd
		var viewerCmd tea.Cmd

		// Update only the focused list
		if m.focusState == focusJournals {
			m.journals, listCmd = m.journals.Update(msg)
			cmds = append(cmds, listCmd)
		} else if m.focusState == focusEntries && m.selectedJournalId != "" {
			// Store previous selected entry ID to detect changes
			previousEntryId := m.selectedEntryId

			m.entries, listCmd = m.entries.Update(msg)
			cmds = append(cmds, listCmd)

			// Check if selection changed and load entry content
			if entry, ok := m.entries.SelectedEntry(); ok {
				if entry.Id != previousEntryId {
					m.selectedEntryId = entry.Id
					cmds = append(cmds, getEntryById(m.ctx, m.database, m.selectedJournalId, m.selectedEntryId))
				}
			} else {
				// No item selected, clear entry
				m.selectedEntryId = ""
				m.viewer.SetContent("")
			}

			// Update textarea while entry list is visible (even if blurred)
			m.viewer, viewerCmd = m.viewer.Update(msg)
			cmds = append(cmds, viewerCmd)
		} else if m.focusState == focusEntry && m.selectedJournalId != "" {
			// Only textarea should react to input in entry focus
			m.viewer, viewerCmd = m.viewer.Update(msg)
			cmds = append(cmds, viewerCmd)
		}
	}

	if m.restoreJournalSelection && m.focusState == focusJournals {
		items := m.journals.Items()
		if len(items) > 0 {
			index := m.lastJournalIndex
			if index < 0 {
				index = 0
			} else if index >= len(items) {
				index = len(items) - 1
			}
			m.journals.Select(index)
		}
		m.restoreJournalSelection = false
	}

	return m, tea.Batch(cmds...)
}

func (m model) entriesViewActive() bool {
	return m.selectedJournalId != ""
}

func (m model) entrySelected() bool {
	return m.selectedEntryId != ""
}

func (m *model) ensureTextareaFocus(active bool) {
	if active {
		if !m.viewer.Focused() {
			m.viewer.Focus()
		}
		return
	}

	if m.viewer.Focused() {
		m.viewer.Blur()
	}
}

func (m *model) setFocusState(next focusState) {
	switch next {
	case focusEntries:
		if !m.entriesViewActive() {
			next = focusJournals
		}
		m.ensureTextareaFocus(false)
	case focusEntry:
		if !m.entriesViewActive() || !m.entrySelected() {
			next = focusEntries
			m.ensureTextareaFocus(false)
		} else {
			m.ensureTextareaFocus(true)
		}
	default:
		next = focusJournals
		m.ensureTextareaFocus(false)
	}

	m.viewerFull = next == focusEntry
	m.focusState = next
}

func (m *model) resizeComponents() {
	if !m.ready {
		return
	}

	if !m.entriesViewActive() {
		m.journals.UpdateWidths(m.width)
		m.journals.SetSize(m.width, m.contentHeight)
		return
	}

	m.entries.UpdateWidths(m.width)
	if m.viewerFull {
		m.entries.SetSize(m.width, 0)
		m.viewer.SetSize(m.width, m.contentHeight)
		return
	}

	listHeight := (m.contentHeight * 60) / 100
	if listHeight < 1 {
		listHeight = 1
	}
	textHeight := m.contentHeight - listHeight
	if textHeight < 1 {
		textHeight = 1
	}

	m.entries.SetSize(m.width, listHeight)
	m.viewer.SetSize(m.width, textHeight)
}

// contentView returns the combined view of both lists (without header/footer)
func (m model) contentView() string {
	// If a journal is selected, show entry list/textarea combo
	if m.entriesViewActive() {
		if m.viewerFull {
			return m.viewer.View()
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			m.entries.View(),
			m.viewer.View(),
		)
	}

	// Otherwise, show only journal list
	return m.journals.View()
}

// Assembles the UI string for each frame
func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	header := m.headerView()
	footer := m.footerView()

	// Build the view: header, content area, footer
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		lipgloss.NewStyle().Width(m.width).Render(m.contentView()),
		footer,
	)
}

func ShowTui(ctx context.Context, database db.Database) error {
	p := tea.NewProgram(initModel(ctx, database), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
