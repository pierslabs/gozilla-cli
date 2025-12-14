package generate

import (
	"github.com/spf13/cobra"
)

var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
	Short:   "Generate code components",
	Long:    `Generate modules, use cases, and other code components for your project`,
}

func init() {
	GenerateCmd.AddCommand(moduleCmd)
}
