//go:build tui

package tui

import (
	"context"
	"fmt"

	"github.com/kompotkot/firn/pkg/db"
	"github.com/kompotkot/firn/pkg/kb"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Keymap for help panel in footer
type keymap = struct {
	enter key.Binding
	esc   key.Binding
	quit  key.Binding
}

func initKeymap() keymap {
	return keymap{
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

type journalsLoadedMsg struct {
	journals []kb.Journal
}

type entriesLoadedMsg struct {
	journalId string
	entries   []kb.Entry
}

type entryLoadedMsg struct {
	journalId string
	entryId   string
	entry     *kb.Entry
}

type errMsg struct {
	operation string
	err       error
}

func (e errMsg) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

// List journals from the database and return as tea data
func listJournals(ctx context.Context, database db.Database, orderByDesc bool, limit int, offset int) tea.Cmd {
	return func() tea.Msg {
		currentCtx := ctx
		if currentCtx == nil {
			currentCtx = context.Background()
		}

		journals, err := database.ListJournals(currentCtx, orderByDesc, limit, offset)
		if err != nil {
			return errMsg{operation: "listJournals", err: err}
		}
		return journalsLoadedMsg{journals: journals}
	}
}

// List journal entries from the database and return as tea data
func listEntries(ctx context.Context, database db.Database, journalId string, orderByDesc bool, limit int, offset int) tea.Cmd {
	return func() tea.Msg {
		currentCtx := ctx
		if currentCtx == nil {
			currentCtx = context.Background()
		}

		entries, err := database.ListEntries(currentCtx, journalId, orderByDesc, limit, offset)
		if err != nil {
			return errMsg{operation: fmt.Sprintf("listEntries(%s)", journalId), err: err}
		}
		return entriesLoadedMsg{journalId: journalId, entries: entries}
	}
}

// Get entry by ID from the database and return as tea data
func getEntryById(ctx context.Context, database db.Database, journalId, entryId string) tea.Cmd {
	return func() tea.Msg {
		currentCtx := ctx
		if currentCtx == nil {
			currentCtx = context.Background()
		}

		entry, err := database.GetEntryById(currentCtx, journalId, entryId)
		if err != nil {
			return errMsg{operation: fmt.Sprintf("getEntryById(%s,%s)", journalId, entryId), err: err}
		}
		return entryLoadedMsg{journalId: journalId, entryId: entryId, entry: entry}
	}
}
