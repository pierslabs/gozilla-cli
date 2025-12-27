package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	templates "github.com/pierslabs/gozilla/internal/templates/module"
)

type ModuleGenerator struct{}

func NewModuleGenerator() *ModuleGenerator {
	return &ModuleGenerator{}
}

func (g *ModuleGenerator) Generate(moduleName string, dependencies []string) error {
	data := templates.ModuleData{
		ModuleName:       moduleName,
		ModuleNameTitle:  strings.Title(moduleName),
		ModuleNameUpper:  strings.ToUpper(moduleName),
		EntityName:       strings.TrimSuffix(moduleName, "s"),
		EntityNamePlural: moduleName,
		Dependencies:     dependencies,
	}

	// Ensure entity name is properly capitalized
	if len(data.EntityName) > 0 {
		data.EntityName = strings.ToUpper(data.EntityName[:1]) + data.EntityName[1:]
	}

	moduleDir := filepath.Join("internal", "modules", moduleName)

	// Create directory structure
	if err := g.createDirectories(moduleDir); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate files
	if err := g.generateFiles(moduleDir, data); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Update container
	containerUpdater := NewContainerUpdater()
	if err := containerUpdater.AddModule(moduleName); err != nil {
		return fmt.Errorf("failed to update container: %w", err)
	}

	return nil
}

func (g *ModuleGenerator) createDirectories(moduleDir string) error {
	dirs := []string{
		moduleDir,
		filepath.Join(moduleDir, "domain"),
		filepath.Join(moduleDir, "application", "dto"),
		filepath.Join(moduleDir, "application", "usecases"),
		filepath.Join(moduleDir, "infra"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (g *ModuleGenerator) generateFiles(moduleDir string, data templates.ModuleData) error {
	files := map[string]string{
		// Module file
		filepath.Join(moduleDir, fmt.Sprintf("%s.module.go", data.ModuleName)): templates.ModuleTemplate(data),

		// Domain layer
		filepath.Join(moduleDir, "domain", fmt.Sprintf("%s.go", data.ModuleName)):            templates.EntityTemplate(data),
		filepath.Join(moduleDir, "domain", fmt.Sprintf("%s_repository.go", data.ModuleName)): templates.RepositoryInterfaceTemplate(data),
		filepath.Join(moduleDir, "domain", "errors.go"):                                      templates.ErrorsTemplate(data),

		// Application layer
		filepath.Join(moduleDir, "application", "dto", fmt.Sprintf("create_%s_dto.go", data.ModuleName)):  templates.CreateDTOTemplate(data),
		filepath.Join(moduleDir, "application", "dto", fmt.Sprintf("update_%s_dto.go", data.ModuleName)):  templates.UpdateDTOTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("create_%s.go", data.ModuleName)): templates.CreateUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("get_%s.go", data.ModuleName)):    templates.GetUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("list_%s.go", data.ModuleName)):   templates.ListUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("update_%s.go", data.ModuleName)): templates.UpdateUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("delete_%s.go", data.ModuleName)): templates.DeleteUseCaseTemplate(data),

		// Infrastructure layer
		filepath.Join(moduleDir, "infra", "handler.go"):                                     templates.HandlerTemplate(data),
		filepath.Join(moduleDir, "infra", "routes.go"):                                      templates.RoutesTemplate(data),
		filepath.Join(moduleDir, "infra", fmt.Sprintf("%s_repository.go", data.ModuleName)): templates.RepositoryImplTemplate(data),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
	}

	return nil
}
