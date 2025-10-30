package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kompotkot/firn/internal/types"
)

// Default configuration values
const (
	DefaultLoggerLevel  = "info"
	DefaultLoggerFormat = "text"

	DefaultDatabaseType            = "sqlite"
	DefaultDatabaseSqliteURI       = "firn.sqlite"
	DefaultDatabasePsqlURI         = "postgres://postgres:postgres@localhost:5432/firn"
	DefaultDatabaseMaxConns        = 10
	DefaultDatabaseConnMaxLifetime = 30 * time.Second
)

// Load and parse configuration
// TODO(kompotkot): Re-write based on https://github.com/kelseyhightower/envconfig
func Load() (*types.Config, error) {
	var cfg types.Config

	logLevelEnv := os.Getenv("LOG_LEVEL")
	if logLevelEnv == "" {
		logLevelEnv = DefaultLoggerLevel
	}

	logFormatEnv := os.Getenv("LOG_FORMAT")
	if logFormatEnv == "" {
		logFormatEnv = DefaultLoggerFormat
	}

	databaseTypeEnv := os.Getenv("DATABASE_TYPE")
	if databaseTypeEnv == "" {
		databaseTypeEnv = DefaultDatabaseType
	}

	databaseURIEnv := os.Getenv("DATABASE_URI")
	if databaseURIEnv == "" {
		switch databaseTypeEnv {
		case "psql":
			databaseURIEnv = DefaultDatabasePsqlURI
		case "sqlite":
			databaseURIEnv = DefaultDatabaseSqliteURI
		default:
			return nil, fmt.Errorf("invalid database type: %s", databaseTypeEnv)
		}
	}

	var databaseMaxConns int
	databaseMaxConnsEnv := os.Getenv("DATABASE_MAX_OPEN_CONNS")
	if databaseMaxConnsEnv != "" {
		if val, err := strconv.Atoi(databaseMaxConnsEnv); err != nil {
			return nil, fmt.Errorf("invalid max open conns: %s, must be a number", databaseMaxConnsEnv)
		} else {
			databaseMaxConns = val
		}
	} else {
		databaseMaxConns = DefaultDatabaseMaxConns
	}

	var databaseConnMaxLifetime time.Duration
	databaseConnMaxLifetimeSecEnv := os.Getenv("DATABASE_CONN_MAX_LIFETIME_SEC")
	if databaseConnMaxLifetimeSecEnv != "" {
		if val, err := strconv.Atoi(databaseConnMaxLifetimeSecEnv); err != nil {
			return nil, fmt.Errorf("invalid conn max lifetime: %s, must be a number", databaseConnMaxLifetimeSecEnv)
		} else {
			databaseConnMaxLifetime = time.Duration(val) * time.Second
		}
	} else {
		databaseConnMaxLifetime = DefaultDatabaseConnMaxLifetime
	}

	cfg = types.Config{
		Logger: types.LoggerConfig{
			Level:  logLevelEnv,
			Format: logFormatEnv,
		},
		Database: types.DatabaseConfig{
			Type:            databaseTypeEnv,
			URI:             databaseURIEnv,
			MaxConns:        databaseMaxConns,
			ConnMaxLifetime: databaseConnMaxLifetime,
		},
	}

	return &cfg, nil
}
