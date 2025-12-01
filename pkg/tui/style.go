//go:build tui

package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Datetime format HH:MM:SS DD-MM-YYYY
	datetimeFormat = "15:04:05 02-01-2006"

	// 3 - from different paddings to not cut datetime at right side
	magicWidthPaddingNum = 3

	// appStyle = lipgloss.NewStyle().Padding(1, 2)

	// --- Header ----
	// Get terminal default colors and invert them
	// Use AdaptiveColor to get colors based on terminal theme
	// headerStyle = invertFBGColors(
	// 	lipgloss.NewStyle().
	// 		Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}). // Default text color
	// 		Background(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#000000"}), // Default background color
	// )
	headerStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false)

	// --- List styles ---

	rightPaddingDatetime = 3 // from different paddings to not cut datetime at right side

	listTitleStyle = lipgloss.NewStyle().UnsetBackground().Bold(true)

	// Style of title of selected journal in list
	listSelectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				Bold(true).
				Padding(0, 0, 0, 1) // Padding from left side

	// Style of description of selected journal in list
	listSelectedDescStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				Padding(0, 0, 0, 1) // Padding from left side

	// Style for NoItems
	listNoItemsStyle = lipgloss.NewStyle()

	// --- Footer ---
	helpStyle  = lipgloss.NewStyle().Align(lipgloss.Left).Padding(0, 0, 0, 1).Faint(true)
	debugStyle = lipgloss.NewStyle().Align(lipgloss.Right).Foreground(lipgloss.Color("#ff0000"))
)

// headerView represents the header view of the TUI
func (m model) headerView() string {
	return headerStyle.Width(m.width).Align(lipgloss.Center).Render("Firn")
}

// footerView represents the footer view of the TUI
func (m model) footerView() string {
	// Show different help based on focus state
	var helpBindings []key.Binding
	switch {
	case m.focusState == focusJournals:
		helpBindings = []key.Binding{m.keys.quit, m.keys.enter}
	case m.focusState == focusEntries && m.selectedJournalId != "":
		helpBindings = []key.Binding{m.keys.esc, m.keys.enter}
	case m.focusState == focusEntry:
		helpBindings = []key.Binding{m.keys.esc}
	default:
		helpBindings = []key.Binding{m.keys.quit}
	}
	help := helpStyle.Render(renderHelpBindings(helpBindings))

	// Compose debug string
	var debugStr string
	if m.debugActive {
		debugStr = fmt.Sprintf("%s [DEBUG]", m.debugStr)
	}

	debug := debugStyle.Width(m.width - len(debugStr)).Render(debugStr)
	return lipgloss.NewStyle().Width(m.width).Render(help + debug)
}

func renderHelpBindings(bindings []key.Binding) string {
	if len(bindings) == 0 {
		return ""
	}

	parts := make([]string, 0, len(bindings))
	for _, b := range bindings {
		if !b.Enabled() {
			continue
		}

		helpInfo := b.Help()
		keyStr := helpInfo.Key
		desc := helpInfo.Desc

		if keyStr == "" {
			keyStr = strings.Join(b.Keys(), "/")
		}

		switch {
		case keyStr != "" && desc != "":
			parts = append(parts, fmt.Sprintf("%s %s", keyStr, desc))
		case keyStr != "":
			parts = append(parts, keyStr)
		case desc != "":
			parts = append(parts, desc)
		}
	}

	return strings.Join(parts, " | ")
}

// invertFBGColors creates a new style with inverted foreground and background colors
// from the original style. All other properties (width, height, padding, margins, etc.)
// are preserved from the original style.
func invertFBGColors(style lipgloss.Style) lipgloss.Style {
	// Get current colors
	fgColor := style.GetForeground()
	bgColor := style.GetBackground()

	// Copy the style (Style is a value type, so assignment copies all properties)
	inverted := style

	inverted = inverted.UnsetForeground()
	inverted = inverted.UnsetBackground()

	// Swap colors
	inverted = inverted.Foreground(bgColor)
	inverted = inverted.Background(fgColor)

	return inverted
}
