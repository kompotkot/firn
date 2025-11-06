//go:build tui

package tui

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/kompotkot/firn/pkg/db"

	tea "github.com/charmbracelet/bubbletea"
)

func initDebug() bool {
	debugActiveRaw := os.Getenv("TUI_DEBUG_ACTIVE")
	debugActive, err := strconv.ParseBool(debugActiveRaw)
	if err != nil {
		fmt.Printf("Error parsing debug variable '%s': %v\n", debugActiveRaw, err)
	}

	return debugActive
}

// List journals from the database and return as tea data
func listJournals(database db.Database, orderByDesc bool, limit int, offset int) tea.Cmd {
	return func() tea.Msg {
		journals, err := database.ListJournals(context.Background(), orderByDesc, limit, offset)
		if err != nil {
			return err
		}
		return journals
	}
}
