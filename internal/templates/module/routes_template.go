package templates

import "fmt"

func RoutesTemplate(data ModuleData) string {
	return fmt.Sprintf(`package infra

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *%sHandler) {
	%s := r.Group("/%s")
	{
		%s.POST("", handler.Create)
		%s.GET("", handler.List)
		%s.GET("/:id", handler.Get)
		%s.PUT("/:id", handler.Update)
		%s.DELETE("/:id", handler.Delete)
	}
}
`, data.ModuleNameTitle,
		data.ModuleName, data.ModuleName,
		data.ModuleName, data.ModuleName, data.ModuleName,
		data.ModuleName, data.ModuleName)
}
