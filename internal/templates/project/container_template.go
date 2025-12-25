package templates

import "fmt"

func ContainerTemplate(data ProjectData) string {
	return fmt.Sprintf(`package container

import (
	"database/sql"

	"%s/internal/modules/health"
	"github.com/gin-gonic/gin"
)

type Container struct {
	DB           *sql.DB
	HealthModule *health.HealthModule
}

func NewContainer(db *sql.DB) *Container {
	return &Container{
		DB:           db,
		HealthModule: health.NewHealthModule(),
	}
}

func (c *Container) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	c.HealthModule.RegisterRoutes(api)
}
`, data.ModulePath)
}
