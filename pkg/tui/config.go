//go:build tui

package tui

import (
	"os"
	"strconv"
)

func initDebug() bool {
	var debugActive bool
	var err error

	debugActiveRaw := os.Getenv("TUI_DEBUG_ACTIVE")
	debugActive, err = strconv.ParseBool(debugActiveRaw)
	if err != nil {
		// TODO(kompotkot): Handle by logger here
		// fmt.Printf("Error parsing debug variable '%s', set to %v, err: %v\n", debugActiveRaw, debugActive, err)
	}

	return debugActive
}
