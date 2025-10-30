//go:build sqlite

package sqlite

import (
	"github.com/kompotkot/firn/pkg/db"
)

func init() {
	// Register SQLite factory when this package is imported
	db.RegisterDatabase(NewFactory())
}
