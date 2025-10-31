package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"os"

	"github.com/kompotkot/firn/internal/config"
	"github.com/kompotkot/firn/internal/logger"
	"github.com/kompotkot/firn/pkg/db"
)

const FIRN_VERSION = "0.0.1"

var (
	moduleTui func(context.Context, db.Database) error
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: firn <command> [args]")
		fmt.Println("Available commands: server, tui, version")
		os.Exit(1)
	}

	var rModule string

	switch os.Args[1] {
	case "server":
		rModule = "server"
	case "tui":
		// TODO(kompotkot): Check if build with tui
		rModule = "tui"
	case "version":
		fmt.Println(FIRN_VERSION)
		return
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

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

	// Create context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch rModule {
	case "server":
		// TODO(kompotkot): Initialize server in go-routine
		fmt.Println("Not implemented yet")
		return
	case "tui":
		log.Info("Starting TUI")
		if err := moduleTui(ctx, database); err != nil {
			log.Error("TUI error", "error", err)
		}
		stop()
	}

	// Wait for shutdown signal
	<-ctx.Done()
	log.Info("Received shutdown signal, starting graceful shutdown")

	// Gracefully close database connection
	log.Info("Closing database connection")
	database.Close()

	log.Info("Application shutdown complete")
}
