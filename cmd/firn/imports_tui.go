//go:build tui

package main

import (
	"context"

	"github.com/kompotkot/firn/pkg/db"

	"github.com/kompotkot/firn/pkg/tui" // Import to register Terminal User Interface
)

func ModuleTui(ctx context.Context, database db.Database) error {
	return tui.ShowTui(ctx, database)
}

func init() {
	moduleTui = ModuleTui
}
