//go:build tui

package tui

import (
	"fmt"

	"github.com/kompotkot/firn/pkg/journal"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
)

// --- Journal Item ---

// jItem represents a journal item in the list
type jItem struct {
	journal journal.Journal

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

func initJournalList() list.Model {
	// Journal list item styles
	jld := list.NewDefaultDelegate()
	jld.Styles.SelectedTitle = listSelectedTitleStyle
	jld.Styles.SelectedDesc = listSelectedDescStyle

	// Initialize journals list (dimensions will be set when window size is received)
	jl := list.New([]list.Item{}, jld, 0, 0)
	jl.Title = "Journals"
	jl.SetFilteringEnabled(false)
	jl.SetShowPagination(true)
	jl.Paginator.Type = paginator.Arabic
	jl.Styles.Title = listTitleStyle
	jl.Styles.NoItems = listNoItemsStyle

	return jl
}
