//go:build tui

package tui

import (
	"fmt"

	"github.com/kompotkot/firn/pkg/kb"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// initList initializes a list of items
func initList(title string) list.Model {
	// List item styles
	ld := list.NewDefaultDelegate()
	ld.Styles.SelectedTitle = listSelectedTitleStyle
	ld.Styles.SelectedDesc = listSelectedDescStyle

	// Initialize list (dimensions will be set when window size is received)
	l := list.New([]list.Item{}, ld, 0, 0)
	l.Title = title
	l.SetFilteringEnabled(false)
	l.SetShowPagination(true)
	l.Paginator.Type = paginator.Arabic
	l.Styles.Title = listTitleStyle
	l.Styles.NoItems = listNoItemsStyle

	return l
}

// Journal item

// jItem represents a journal item in the list
type jItem struct {
	journal kb.Journal

	widthTitle int // Width for Title
	widthDesc  int // Width for Description
}

func (i jItem) Title() string {
	name := lipgloss.NewStyle().Render(i.journal.Name)

	// Create right-aligned updated_at with widthTitle
	updatedAt := lipgloss.NewStyle().Width(i.widthTitle).Align(lipgloss.Right).Render(i.journal.UpdatedAt.Format(datetimeFormat))

	return name + updatedAt
}

func (i jItem) Description() string {
	id := fmt.Sprintf("ID: %s", i.journal.Id)

	// Create right-aligned "Created At: {timestamp}" with widthDesc
	createdAtText := fmt.Sprintf("Created At: %s", i.journal.CreatedAt.Format(datetimeFormat))
	createdAt := lipgloss.NewStyle().Width(i.widthDesc).Align(lipgloss.Right).Render(createdAtText)

	return id + createdAt
}

func (i jItem) FilterValue() string { return i.journal.Name }

// Entry item

// eItem represents journal entry item in the list
type eItem struct {
	entry kb.Entry

	widthTitle int // Width for Title
	widthDesc  int // Width for Description
}

func (i eItem) Title() string {
	title := lipgloss.NewStyle().Render(i.entry.Title)
	width := i.widthTitle
	if width < 0 {
		width = 0
	}
	updatedAt := lipgloss.NewStyle().Width(width).Align(lipgloss.Right).Render(i.entry.UpdatedAt.Format(datetimeFormat))

	return title + updatedAt
}

func (i eItem) Description() string {
	id := fmt.Sprintf("ID: %s", i.entry.Id)

	// Create right-aligned "Created At: {timestamp}" with widthDesc
	createdAtText := fmt.Sprintf("Created At: %s", i.entry.CreatedAt.Format(datetimeFormat))
	createdAt := lipgloss.NewStyle().Width(i.widthDesc).Align(lipgloss.Right).Render(createdAtText)

	return id + createdAt
}

func (i eItem) FilterValue() string { return i.entry.Title }

// initTextarea initializes an entry textarea
func initTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "No entry selected"
	ta.CharLimit = 0 // No limit
	ta.SetWidth(80)  // Will be updated on window resize
	ta.SetHeight(5)  // Will be updated on window resize
	ta.ShowLineNumbers = false
	ta.Blur() // Blur textarea so it doesn't intercept keyboard input

	return ta
}
