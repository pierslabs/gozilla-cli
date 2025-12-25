package templates

import "fmt"

func HealthModuleTemplate(data ProjectData) string {
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
