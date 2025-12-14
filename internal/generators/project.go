package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ProjectGenerator struct{}

func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

type ProjectData struct {
	ProjectName string
	ModulePath  string
	ProjectDir  string
}

func (g *ProjectGenerator) Generate(projectName, projectDir string) error {
	data := ProjectData{
		ProjectName: projectDir,
		ModulePath:  projectName,
		ProjectDir:  projectDir,
	}

	// Create directory structure
	if err := g.createDirectories(projectDir); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files
	if err := g.generateFiles(data); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Initialize go module
	if err := g.initGoModule(projectDir, projectName); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	return nil
}

func (g *ProjectGenerator) createDirectories(projectDir string) error {
	dirs := []string{
		projectDir,
		filepath.Join(projectDir, "cmd", "api"),
		filepath.Join(projectDir, "internal", "domain"),
		filepath.Join(projectDir, "internal", "infrastructure", "database"),
		filepath.Join(projectDir, "internal", "infrastructure", "http"),
		filepath.Join(projectDir, "internal", "infrastructure", "config"),
		filepath.Join(projectDir, "internal", "infrastructure", "container"),
		filepath.Join(projectDir, "internal", "modules", "health"),
		filepath.Join(projectDir, "internal", "modules", "health", "domain"),
		filepath.Join(projectDir, "internal", "modules", "health", "application"),
		filepath.Join(projectDir, "internal", "modules", "health", "infra"),
		filepath.Join(projectDir, "pkg"),
		filepath.Join(projectDir, "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (g *ProjectGenerator) generateFiles(data ProjectData) error {
	files := map[string]string{
		filepath.Join(data.ProjectDir, "cmd", "api", "main.go"):                                          g.mainGoTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "config", "config.go"):              g.configTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "database", "database.go"):          g.databaseTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "http", "server.go"):                g.serverTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "container", "container.go"):        g.containerTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "health.module.go"):              g.healthModuleTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "infra", "handler.go"):           g.healthHandlerTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "infra", "routes.go"):            g.healthRoutesTemplate(data),
		filepath.Join(data.ProjectDir, "docker-compose.yaml"):                                            g.dockerComposeTemplate(data),
		filepath.Join(data.ProjectDir, "Makefile"):                                                       g.makefileTemplate(data),
		filepath.Join(data.ProjectDir, ".env.example"):                                                   g.envExampleTemplate(data),
		filepath.Join(data.ProjectDir, ".gitignore"):                                                     g.gitignoreTemplate(data),
		filepath.Join(data.ProjectDir, "README.md"):                                                      g.readmeTemplate(data),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
	}

	return nil
}

func (g *ProjectGenerator) initGoModule(projectDir, modulePath string) error {
	goModPath := filepath.Join(projectDir, "go.mod")

	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)
`, modulePath)

	return os.WriteFile(goModPath, []byte(content), 0644)
}

// Templates

func (g *ProjectGenerator) mainGoTemplate(data ProjectData) string {
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

func (g *ProjectGenerator) configTemplate(data ProjectData) string {
	return `package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/` + data.ProjectName + `?sslmode=disable"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
}

func (g *ProjectGenerator) databaseTemplate(data ProjectData) string {
	return `package database

import (
	"database/sql"
	"fmt"
)

func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
`
}

func (g *ProjectGenerator) serverTemplate(data ProjectData) string {
	return fmt.Sprintf(`package http

import (
	"fmt"

	"%s/internal/infrastructure/config"
	"%s/internal/infrastructure/container"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router    *gin.Engine
	config    *config.Config
	container *container.Container
}

func NewServer(cfg *config.Config, c *container.Container) *Server {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	s := &Server{
		router:    router,
		config:    cfg,
		container: c,
	}

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// Register module routes
	s.container.RegisterRoutes(s.router)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%%s", s.config.Port)
	fmt.Printf("Server starting on %%s\n", addr)
	return s.router.Run(addr)
}
`, data.ModulePath, data.ModulePath)
}

func (g *ProjectGenerator) containerTemplate(data ProjectData) string {
	return fmt.Sprintf(`package container

import (
	"database/sql"

	"%s/internal/modules/health"
	"github.com/gin-gonic/gin"
)

type Container struct {
	DB           *sql.DB
	HealthModule *health.HealthModule
}

func NewContainer(db *sql.DB) *Container {
	return &Container{
		DB:           db,
		HealthModule: health.NewHealthModule(),
	}
}

func (c *Container) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	c.HealthModule.RegisterRoutes(api)
}
`, data.ModulePath)
}

func (g *ProjectGenerator) healthModuleTemplate(data ProjectData) string {
	return fmt.Sprintf(`package health

import (
	"%s/internal/modules/health/infra"
	"github.com/gin-gonic/gin"
)

type HealthModule struct {
	Handler *infra.HealthHandler
}

func NewHealthModule() *HealthModule {
	handler := infra.NewHealthHandler()
	return &HealthModule{
		Handler: handler,
	}
}

func (m *HealthModule) RegisterRoutes(r *gin.RouterGroup) {
	infra.RegisterRoutes(r, m.Handler)
}
`, data.ModulePath)
}

func (g *ProjectGenerator) healthHandlerTemplate(data ProjectData) string {
	return `package infra

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"message": "Service is healthy",
	})
}
`
}

func (g *ProjectGenerator) healthRoutesTemplate(data ProjectData) string {
	return `package infra

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *HealthHandler) {
	health := r.Group("/health")
	{
		health.GET("", handler.Check)
	}
}
`
}

func (g *ProjectGenerator) dockerComposeTemplate(data ProjectData) string {
	return fmt.Sprintf(`version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: %s
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
`, data.ProjectName)
}

func (g *ProjectGenerator) makefileTemplate(data ProjectData) string {
	return `.PHONY: run build test docker-up docker-down migrate-up migrate-down

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	@echo "Migrations not configured yet"

migrate-down:
	@echo "Migrations not configured yet"

clean:
	rm -rf bin/
`
}

func (g *ProjectGenerator) envExampleTemplate(data ProjectData) string {
	return fmt.Sprintf(`PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/%s?sslmode=disable
ENVIRONMENT=development
`, data.ProjectName)
}

func (g *ProjectGenerator) gitignoreTemplate(data ProjectData) string {
	return `# Binaries
bin/
*.exe
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment files
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Database
*.db
*.sqlite
`
}

func (g *ProjectGenerator) readmeTemplate(data ProjectData) string {
	projectName := strings.Title(data.ProjectName)
	return fmt.Sprintf(`# %s

Generated with [Gozilla](https://github.com/pierslabs/gozilla)

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Installation

1. Start the database:
`+"```bash"+`
make docker-up
`+"```"+`

2. Run the application:
`+"```bash"+`
make run
`+"```"+`

The API will be available at http://localhost:8080

### API Endpoints

- `+"`GET /api/v1/health`"+` - Health check

## Development

### Generate a new module

`+"```bash"+`
gozilla generate module users
`+"```"+`

### Run tests

`+"```bash"+`
make test
`+"```"+`

### Build

`+"```bash"+`
make build
`+"```"+`

## Project Structure

`+"```"+`
.
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                # Shared domain
│   ├── infrastructure/        # Infrastructure layer
│   │   ├── config/           # Configuration
│   │   ├── database/         # Database connection
│   │   ├── http/             # HTTP server
│   │   └── container/        # DI container
│   └── modules/              # Feature modules
│       └── health/           # Health check module
├── migrations/               # Database migrations
└── docker-compose.yaml       # Docker setup
`+"```"+`

## License

MIT
`, projectName)
}
