//go:build psql

package psql

import (
	"github.com/kompotkot/firn/pkg/db"
)

func init() {
	// Register PostgreSQL factory when this package is imported
	db.RegisterDatabase(NewFactory())
}
