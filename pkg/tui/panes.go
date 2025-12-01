//go:build tui

package tui

import (
	"fmt"

	"github.com/kompotkot/firn/pkg/kb"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type journalsPane struct {
	list list.Model
}

func newJournalsPane() journalsPane {
	return journalsPane{
		list: initList("Journals"),
	}
}

func (p journalsPane) Update(msg tea.Msg) (journalsPane, tea.Cmd) {
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p journalsPane) View() string {
	return p.list.View()
}

func (p *journalsPane) SetItems(items []list.Item) {
	p.list.SetItems(items)
}

func (p *journalsPane) SetTotalPages(total int) {
	p.list.Paginator.SetTotalPages(total)
}

func (p *journalsPane) SetSize(width, height int) {
	p.list.SetSize(width, height)
}

func (p journalsPane) SelectedJournal() (kb.Journal, bool) {
	item := p.list.SelectedItem()
	if item == nil {
		return kb.Journal{}, false
	}
	if ji, ok := item.(jItem); ok {
		return ji.journal, true
	}
	return kb.Journal{}, false
}

func (p journalsPane) SelectedItem() list.Item {
	return p.list.SelectedItem()
}

func (p journalsPane) GlobalIndex() int {
	return p.list.GlobalIndex()
}

func (p *journalsPane) Select(index int) {
	p.list.Select(index)
}

func (p journalsPane) Items() []list.Item {
	return p.list.Items()
}

func (p *journalsPane) UpdateWidths(width int) {
	currentItems := p.list.Items()
	updated := make([]list.Item, len(currentItems))
	for i, it := range currentItems {
		if ji, ok := it.(jItem); ok {
			ji.widthTitle = width - len(ji.journal.Name) - rightPaddingDatetime
			ji.widthDesc = width - len(fmt.Sprintf("ID: %s", ji.journal.Id)) - rightPaddingDatetime
			updated[i] = ji
		} else {
			updated[i] = it
		}
	}
	p.list.SetItems(updated)
}

type entriesPane struct {
	list list.Model
}

func newEntriesPane() entriesPane {
	return entriesPane{
		list: initList("Entries"),
	}
}

func (p entriesPane) Update(msg tea.Msg) (entriesPane, tea.Cmd) {
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p entriesPane) View() string {
	return p.list.View()
}

func (p *entriesPane) SetItems(items []list.Item) {
	p.list.SetItems(items)
}

func (p *entriesPane) SetTotalPages(total int) {
	p.list.Paginator.SetTotalPages(total)
}

func (p *entriesPane) SetSize(width, height int) {
	p.list.SetSize(width, height)
}

func (p *entriesPane) SetTitle(title string) {
	p.list.Title = title
}

func (p entriesPane) SelectedItem() list.Item {
	return p.list.SelectedItem()
}

func (p entriesPane) SelectedEntry() (kb.Entry, bool) {
	item := p.list.SelectedItem()
	if item == nil {
		return kb.Entry{}, false
	}
	if ei, ok := item.(eItem); ok {
		return ei.entry, true
	}
	return kb.Entry{}, false
}

func (p entriesPane) Items() []list.Item {
	return p.list.Items()
}

func (p *entriesPane) UpdateWidths(width int) {
	currentItems := p.list.Items()
	updated := make([]list.Item, len(currentItems))
	for i, it := range currentItems {
		if ei, ok := it.(eItem); ok {
			idLabel := fmt.Sprintf("ID: %s", ei.entry.Id)
			ei.widthTitle = width - len(ei.entry.Title) - rightPaddingDatetime
			ei.widthDesc = width - len(idLabel) - rightPaddingDatetime
			updated[i] = ei
		} else {
			updated[i] = it
		}
	}
	p.list.SetItems(updated)
}

func (p *entriesPane) Clear() {
	p.list.SetItems([]list.Item{})
	p.list.Paginator.SetTotalPages(0)
}

// Entry viewer textarea

type entryViewer struct {
	textarea textarea.Model
}

func newEntryViewer() entryViewer {
	return entryViewer{
		textarea: initTextarea(),
	}
}

func (v entryViewer) Update(msg tea.Msg) (entryViewer, tea.Cmd) {
	var cmd tea.Cmd
	v.textarea, cmd = v.textarea.Update(msg)
	return v, cmd
}

func (v entryViewer) View() string {
	return v.textarea.View()
}

func (v *entryViewer) SetSize(width, height int) {
	v.textarea.SetWidth(width)
	v.textarea.SetHeight(height)
}

func (v *entryViewer) SetContent(value string) {
	v.textarea.SetValue(value)
}

func (v entryViewer) Focused() bool {
	return v.textarea.Focused()
}

func (v *entryViewer) Focus() {
	v.textarea.Focus()
}

func (v *entryViewer) Blur() {
	v.textarea.Blur()
}
