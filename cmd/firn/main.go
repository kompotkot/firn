package firn

import (
	"context"
	"fmt"

	"os"

	"github.com/kompotkot/firn/internal/config"
	"github.com/kompotkot/firn/internal/logger"
	"github.com/kompotkot/firn/pkg/db"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logger)
	log.Info("Logger initialized")

	// Initialize database connection using registry
	log.Info("Initializing database connection")
	database, err := db.CreateDatabase(
		cfg.Database.Type,
		cfg.Database.URI,
		cfg.Database.MaxConns,
		int64(cfg.Database.ConnMaxLifetime),
	)
	if err != nil {
		log.Error("Failed to initialize database connection", "error", err)
		os.Exit(1)
	}

	// Test database connection
	if err := database.TestConnection(context.Background()); err != nil {
		log.Error("Failed to test database connection", "error", err)
		os.Exit(1)
	}
	log.Info("Database connection established successfully")

	// Gracefully close database connection
	log.Info("Closing database connection")
	database.Close()

	log.Info("Application shutdown complete")
}
