//go:build tui

package tui

import "github.com/charmbracelet/lipgloss"


var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	
	listSelectedStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(lipgloss.Color("62")).
	Foreground(lipgloss.Color("62")).
	Padding(0, 0, 0, 1)
)

// headerView represents the header view of the TUI
func (m model) headerView() string {
	title := "Firn"
	return lipgloss.JoinHorizontal(lipgloss.Center, title)
}
