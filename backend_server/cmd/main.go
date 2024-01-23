package main

import (
	"backend_server/configs"
	"backend_server/internal/server"
	"backend_server/pkg/connection/postgres"
	"context"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Printf("LoadConfig: %w", err)
	}

	dbManager, err := postgres.NewDatabase(ctx, cfg)
	if err != nil {
		log.Fatalf("Could not connect to PostgresDB:", err)
	}

	defer func(database *postgres.Database) {
		err := database.Close()
		if err != nil {
			log.Printf("Error closing Postgres database: %s", err.Error())
		}
	}(dbManager)

	app := server.NewServer(cfg, dbManager)

	if err := app.Run(); err != nil {
		log.Fatalf("Cannot start server: %w", err)
	}
}
