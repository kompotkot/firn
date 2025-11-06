//go:build tui

package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

// Keymap for help panel in footer
type keymap = struct {
	quit key.Binding
}

func initKeymap() keymap {
	return keymap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
