package templates

import (
	"os"
	"strings"
)

type ModuleData struct {
	ModuleName       string
	ModuleNameTitle  string
	ModuleNameUpper  string
	EntityName       string
	EntityNamePlural string
	Dependencies     []string
}

// GetModulePath returns the module path from go.mod
func GetModulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "app" // fallback
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return "app" // fallback
}
