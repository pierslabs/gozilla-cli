package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierslabs/gozilla/internal/generators"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Generate a new Go project with Clean Architecture",
	Long: `Creates a new Go project with:
- Clean Architecture structure
- Gin HTTP framework
- PostgreSQL database setup
- Docker Compose configuration
- DI container skeleton
- Example health check module`,
	Args: cobra.ExactArgs(1),
	Example: `  gozilla new my-api
  gozilla new github.com/myuser/my-project`,
	RunE: runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Clean and validate project name
	projectName = strings.TrimSpace(projectName)
	if strings.Contains(projectName, " ") {
		return fmt.Errorf("project name cannot contain spaces")
	}

	// Extract project directory name from full path if provided
	projectDir := filepath.Base(projectName)

	// Check if directory already exists
	if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectDir)
	}

	fmt.Printf("🚀 Creating new project: %s\n", projectName)

	// Generate project
	generator := generators.NewProjectGenerator()
	if err := generator.Generate(projectName, projectDir); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("\n✅ Project created successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", projectDir)
	fmt.Printf("  docker-compose up -d\n")
	fmt.Printf("  make run\n\n")
	fmt.Printf("Generate your first module:\n")
	fmt.Printf("  gozilla generate module users\n")

	return nil
}
