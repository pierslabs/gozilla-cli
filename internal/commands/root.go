package commands

import (
	"fmt"
	"os"

	"github.com/pierslabs/gozilla/internal/commands/generate"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gozilla",
	Short: "Generate production-ready Go projects with Clean Architecture",
	Long: `Gozilla is a CLI tool that generates Go API projects with:
- Clean Architecture (Domain/Application/Infrastructure layers)
- Automatic dependency injection (via generated code)
- Modular structure (each feature is a self-contained module)
- Production-ready setup (Docker, PostgreSQL, tests, migrations)`,
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
}
