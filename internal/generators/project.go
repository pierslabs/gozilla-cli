package generators

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	templates "github.com/pierslabs/gozilla/internal/templates/project"
)

type ProjectGenerator struct{}

func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

func (g *ProjectGenerator) Generate(projectName, projectDir string) error {
	data := templates.ProjectData{
		ProjectName: projectName,
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
	paths := []string{
		".",
		"cmd/api",
		"internal/domain",
		"internal/infrastructure/database",
		"internal/infrastructure/http",
		"internal/infrastructure/config",
		"internal/infrastructure/container",
		"internal/modules/health/domain",
		"internal/modules/health/application",
		"internal/modules/health/infra",
		"pkg",
		"migrations",
	}

	for _, path := range paths {
		dir := filepath.Join(projectDir, path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (g *ProjectGenerator) generateFiles(data templates.ProjectData) error {
	type fileSpec struct {
		path     string
		template func(templates.ProjectData) string
	}

	files := []fileSpec{
		{"cmd/api/main.go", templates.MainGoTemplate},
		{"internal/infrastructure/config/config.go", templates.ConfigTemplate},
		{"internal/infrastructure/database/database.go", templates.DatabaseTemplate},
		{"internal/infrastructure/http/server.go", templates.ServerTemplate},
		{"internal/infrastructure/container/container.go", templates.ContainerTemplate},
		{"internal/modules/health/health.module.go", templates.HealthModuleTemplate},
		{"internal/modules/health/infra/handler.go", templates.HealthHandlerTemplate},
		{"internal/modules/health/infra/routes.go", templates.HealthRoutesTemplate},
		{"docker-compose.yaml", templates.DockerComposeTemplate},
		{"Makefile", templates.MakefileTemplate},
		{".env.example", templates.EnvExampleTemplate},
		{".gitignore", templates.GitignoreTemplate},
		{"README.md", templates.ReadmeTemplate},
	}

	for _, file := range files {
		fullPath := filepath.Join(data.ProjectDir, file.path)
		content := file.template(data)

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", file.path, err)
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

func (g *ProjectGenerator) InstallDependencies(projectDir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w\n%s", err, string(output))
	}
	return nil
}
