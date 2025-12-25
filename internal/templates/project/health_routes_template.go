package templates

func HealthRoutesTemplate(data ProjectData) string {
	return `package infra

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *HealthHandler) {
	health := r.Group("/health")
	{
		health.GET("", handler.Check)
	}
}
`
}
