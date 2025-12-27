package templates

import "fmt"

func ModuleTemplate(data ModuleData) string {
	return fmt.Sprintf(`package %s

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"%s/internal/modules/%s/application/usecases"
	"%s/internal/modules/%s/infra"
)

type %sModule struct {
	Handler *infra.%sHandler
}

func New%sModule(db *sql.DB) *%sModule {
	repo := infra.New%sRepository(db)

	createUC := usecases.NewCreate%sUseCase(repo)
	getUC := usecases.NewGet%sUseCase(repo)
	listUC := usecases.NewList%sUseCase(repo)
	updateUC := usecases.NewUpdate%sUseCase(repo)
	deleteUC := usecases.NewDelete%sUseCase(repo)

	handler := infra.New%sHandler(createUC, getUC, listUC, updateUC, deleteUC)

	return &%sModule{
		Handler: handler,
	}
}

func (m *%sModule) RegisterRoutes(r *gin.RouterGroup) {
	infra.RegisterRoutes(r, m.Handler)
}
`, data.ModuleName,
		GetModulePath(), data.ModuleName,
		GetModulePath(), data.ModuleName,
		data.ModuleNameTitle, data.ModuleNameTitle,
		data.ModuleNameTitle, data.ModuleNameTitle,
		data.EntityName,
		data.EntityName, data.EntityName, data.ModuleNameTitle,
		data.EntityName, data.EntityName,
		data.ModuleNameTitle,
		data.ModuleNameTitle,
		data.ModuleNameTitle)
}
