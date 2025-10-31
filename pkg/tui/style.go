//go:build tui

package tui

import "github.com/charmbracelet/lipgloss"


// headerView represents the header view of the TUI
func (m model) headerView() string {
	title := "Firn"
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}
