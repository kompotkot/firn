//go:build sqlite

package main

import (
	_ "github.com/kompotkot/firn/pkg/db/sqlite" // Import to register SQLite factory
)
