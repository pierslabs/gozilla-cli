package generators

import (
	"fmt"
	"os"
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
		ModulePath:  fmt.Sprintf("github.com/yourusername/%s", projectName),
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

func (g *ProjectGenerator) generateFiles(data templates.ProjectData) error {
	files := map[string]string{
		filepath.Join(data.ProjectDir, "cmd", "api", "main.go"):                                   templates.MainGoTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "config", "config.go"):       templates.ConfigTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "database", "database.go"):   templates.DatabaseTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "http", "server.go"):         templates.ServerTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "infrastructure", "container", "container.go"): templates.ContainerTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "health.module.go"):       templates.HealthModuleTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "infra", "handler.go"):    templates.HealthHandlerTemplate(data),
		filepath.Join(data.ProjectDir, "internal", "modules", "health", "infra", "routes.go"):     templates.HealthRoutesTemplate(data),
		filepath.Join(data.ProjectDir, "docker-compose.yaml"):                                     templates.DockerComposeTemplate(data),
		filepath.Join(data.ProjectDir, "Makefile"):                                                templates.MakefileTemplate(data),
		filepath.Join(data.ProjectDir, ".env.example"):                                            templates.EnvExampleTemplate(data),
		filepath.Join(data.ProjectDir, ".gitignore"):                                              templates.GitignoreTemplate(data),
		filepath.Join(data.ProjectDir, "README.md"):                                               templates.ReadmeTemplate(data),
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
