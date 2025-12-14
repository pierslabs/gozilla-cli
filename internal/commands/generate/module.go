package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierslabs/gozilla/internal/generators"
	"github.com/spf13/cobra"
)

var (
	moduleDependencies []string
)

var moduleCmd = &cobra.Command{
	Use:     "module [name]",
	Aliases: []string{"mod", "m"},
	Short:   "Generate a new feature module",
	Long: `Generates a complete feature module with:
- Domain layer (entity, repository interface, errors)
- Application layer (DTOs, use cases)
- Infrastructure layer (handlers, repository impl, routes)
- Tests for each layer
- Module DI file
- Auto-updates container.go`,
	Args: cobra.ExactArgs(1),
	Example: `  gozilla generate module users
  gozilla g mod orders --depends=users
  gozilla g m products --depends=users,categories`,
	RunE: runGenerateModule,
}

func init() {
	moduleCmd.Flags().StringSliceVar(&moduleDependencies, "depends", []string{}, "Module dependencies (comma-separated)")
}

func runGenerateModule(cmd *cobra.Command, args []string) error {
	moduleName := args[0]

	// Validate module name
	if moduleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	// Clean and validate module name
	moduleName = strings.TrimSpace(moduleName)
	moduleName = strings.ToLower(moduleName)

	if strings.Contains(moduleName, " ") {
		return fmt.Errorf("module name cannot contain spaces")
	}

	// Check if we're in a gozilla project
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("not in a Go project directory (go.mod not found)")
	}

	modulesDir := filepath.Join("internal", "modules")
	if _, err := os.Stat(modulesDir); os.IsNotExist(err) {
		return fmt.Errorf("not in a gozilla project (internal/modules not found)")
	}

	// Check if module already exists
	moduleDir := filepath.Join(modulesDir, moduleName)
	if _, err := os.Stat(moduleDir); !os.IsNotExist(err) {
		return fmt.Errorf("module '%s' already exists", moduleName)
	}

	fmt.Printf("🔧 Generating module: %s\n", moduleName)

	if len(moduleDependencies) > 0 {
		fmt.Printf("📦 Dependencies: %s\n", strings.Join(moduleDependencies, ", "))
	}

	// Generate module
	generator := generators.NewModuleGenerator()
	if err := generator.Generate(moduleName, moduleDependencies); err != nil {
		return fmt.Errorf("failed to generate module: %w", err)
	}

	fmt.Printf("\n✅ Module '%s' created successfully!\n\n", moduleName)
	fmt.Printf("Generated files:\n")
	fmt.Printf("  internal/modules/%s/\n", moduleName)
	fmt.Printf("    ├── %s.module.go\n", moduleName)
	fmt.Printf("    ├── domain/\n")
	fmt.Printf("    ├── application/\n")
	fmt.Printf("    └── infra/\n\n")
	fmt.Printf("Container updated:\n")
	fmt.Printf("  internal/infrastructure/container/container.go\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("  1. Implement business logic in domain/\n")
	fmt.Printf("  2. Add use cases in application/\n")
	fmt.Printf("  3. Run: make run\n")

	return nil
}
