package templates

import "fmt"

func ServerTemplate(data ProjectData) string {
	return fmt.Sprintf(`package http

import (
	"fmt"

	"%s/internal/infrastructure/config"
	"%s/internal/infrastructure/container"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router    *gin.Engine
	config    *config.Config
	container *container.Container
}

func NewServer(cfg *config.Config, c *container.Container) *Server {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	s := &Server{
		router:    router,
		config:    cfg,
		container: c,
	}

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// Register module routes
	s.container.RegisterRoutes(s.router)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%%s", s.config.Port)
	fmt.Printf("Server starting on %%s\n", addr)
	return s.router.Run(addr)
}
`, data.ModulePath, data.ModulePath)
}
