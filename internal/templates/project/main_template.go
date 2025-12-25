package templates

import "fmt"

type ProjectData struct {
	ProjectName string
	ModulePath  string
	ProjectDir  string
}

func MainGoTemplate(data ProjectData) string {
	return fmt.Sprintf(`package main

import (
    "log"

    "%s/internal/infrastructure/config"
    "%s/internal/infrastructure/container"
    "%s/internal/infrastructure/database"
    "%s/internal/infrastructure/http"
    _ "github.com/lib/pq"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %%v", err)
    }

    // Initialize database
    db, err := database.NewConnection(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %%v", err)
    }
    defer db.Close()

    // Initialize DI container
    container := container.NewContainer(db)

    // Start HTTP server
    server := http.NewServer(cfg, container)
    if err := server.Start(); err != nil {
        log.Fatalf("Failed to start server: %%v", err)
    }
}
`, data.ModulePath, data.ModulePath, data.ModulePath, data.ModulePath)
}
