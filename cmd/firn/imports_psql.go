//go:build psql

package main

import (
	_ "github.com/kompotkot/firn/pkg/db/psql" // Import to register PostgreSQL factory
)
