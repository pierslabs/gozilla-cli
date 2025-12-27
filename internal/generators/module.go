package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ModuleGenerator struct{}

func NewModuleGenerator() *ModuleGenerator {
	return &ModuleGenerator{}
}

type ModuleData struct {
	ModuleName       string
	ModuleNameTitle  string
	ModuleNameUpper  string
	EntityName       string
	EntityNamePlural string
	Dependencies     []string
}

func (g *ModuleGenerator) Generate(moduleName string, dependencies []string) error {
	data := ModuleData{
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

func (g *ModuleGenerator) generateFiles(moduleDir string, data ModuleData) error {
	files := map[string]string{
		// Module file
		filepath.Join(moduleDir, fmt.Sprintf("%s.module.go", data.ModuleName)): g.moduleTemplate(data),

		// Domain layer
		filepath.Join(moduleDir, "domain", fmt.Sprintf("%s.go", data.ModuleName)):            g.entityTemplate(data),
		filepath.Join(moduleDir, "domain", fmt.Sprintf("%s_repository.go", data.ModuleName)): g.repositoryInterfaceTemplate(data),
		filepath.Join(moduleDir, "domain", "errors.go"):                                      g.errorsTemplate(data),

		// Application layer
		filepath.Join(moduleDir, "application", "dto", fmt.Sprintf("create_%s_dto.go", data.ModuleName)):  g.createDTOTemplate(data),
		filepath.Join(moduleDir, "application", "dto", fmt.Sprintf("update_%s_dto.go", data.ModuleName)):  g.updateDTOTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("create_%s.go", data.ModuleName)): g.createUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("get_%s.go", data.ModuleName)):    g.getUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("list_%s.go", data.ModuleName)):   g.listUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("update_%s.go", data.ModuleName)): g.updateUseCaseTemplate(data),
		filepath.Join(moduleDir, "application", "usecases", fmt.Sprintf("delete_%s.go", data.ModuleName)): g.deleteUseCaseTemplate(data),

		// Infrastructure layer
		filepath.Join(moduleDir, "infra", "handler.go"):                                     g.handlerTemplate(data),
		filepath.Join(moduleDir, "infra", "routes.go"):                                      g.routesTemplate(data),
		filepath.Join(moduleDir, "infra", fmt.Sprintf("%s_repository.go", data.ModuleName)): g.repositoryImplTemplate(data),
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}
	}

	return nil
}

// Templates

func (g *ModuleGenerator) moduleTemplate(data ModuleData) string {
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
		getModulePath(), data.ModuleName,
		getModulePath(), data.ModuleName,
		data.ModuleNameTitle, data.ModuleNameTitle,
		data.ModuleNameTitle, data.ModuleNameTitle,
		data.ModuleNameTitle,
		data.EntityName, data.EntityName, data.ModuleNameTitle,
		data.EntityName, data.EntityName,
		data.ModuleNameTitle,
		data.ModuleNameTitle,
		data.ModuleNameTitle)
}

func (g *ModuleGenerator) entityTemplate(data ModuleData) string {
	return fmt.Sprintf(`package domain

import "time"

type %s struct {
	ID        int64     `+"`json:\"id\"`"+`
	Name      string    `+"`json:\"name\"`"+`
	CreatedAt time.Time `+"`json:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\"`"+`
}
`, data.EntityName)
}

func (g *ModuleGenerator) repositoryInterfaceTemplate(data ModuleData) string {
	return fmt.Sprintf(`package domain

import "context"

type %sRepository interface {
	Create(ctx context.Context, %s *%s) error
	GetByID(ctx context.Context, id int64) (*%s, error)
	List(ctx context.Context) ([]*%s, error)
	Update(ctx context.Context, %s *%s) error
	Delete(ctx context.Context, id int64) error
}
`, data.EntityName,
		strings.ToLower(data.EntityName[:1]), data.EntityName,
		data.EntityName, data.EntityName,
		strings.ToLower(data.EntityName[:1]), data.EntityName)
}

func (g *ModuleGenerator) errorsTemplate(data ModuleData) string {
	return fmt.Sprintf(`package domain

import "errors"

var (
	Err%sNotFound = errors.New("%s not found")
	Err%sAlreadyExists = errors.New("%s already exists")
	ErrInvalid%s = errors.New("invalid %s data")
)
`, data.EntityName, strings.ToLower(data.ModuleName),
		data.EntityName, strings.ToLower(data.ModuleName),
		data.EntityName, strings.ToLower(data.ModuleName))
}

func (g *ModuleGenerator) createDTOTemplate(data ModuleData) string {
	return fmt.Sprintf(`package dto

type Create%sDTO struct {
	Name string `+"`json:\"name\" binding:\"required\"`"+`
}
`, data.EntityName)
}

func (g *ModuleGenerator) updateDTOTemplate(data ModuleData) string {
	return fmt.Sprintf(`package dto

type Update%sDTO struct {
	Name string `+"`json:\"name\"`"+`
}
`, data.EntityName)
}

func (g *ModuleGenerator) createUseCaseTemplate(data ModuleData) string {
	entityVar := strings.ToLower(data.EntityName[:1])
	return fmt.Sprintf(`package usecases

import (
	"context"
	"time"

	"%s/internal/modules/%s/application/dto"
	"%s/internal/modules/%s/domain"
)

type Create%sUseCase struct {
	repo domain.%sRepository
}

func NewCreate%sUseCase(repo domain.%sRepository) *Create%sUseCase {
	return &Create%sUseCase{
		repo: repo,
	}
}

func (uc *Create%sUseCase) Execute(ctx context.Context, input dto.Create%sDTO) (*domain.%s, error) {
	%s := &domain.%s{
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, %s); err != nil {
		return nil, err
	}

	return %s, nil
}
`, getModulePath(), data.ModuleName, getModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		entityVar, data.EntityName,
		entityVar,
		entityVar)
}

func (g *ModuleGenerator) getUseCaseTemplate(data ModuleData) string {
	return fmt.Sprintf(`package usecases

import (
	"context"

	"%s/internal/modules/%s/domain"
)

type Get%sUseCase struct {
	repo domain.%sRepository
}

func NewGet%sUseCase(repo domain.%sRepository) *Get%sUseCase {
	return &Get%sUseCase{
		repo: repo,
	}
}

func (uc *Get%sUseCase) Execute(ctx context.Context, id int64) (*domain.%s, error) {
	return uc.repo.GetByID(ctx, id)
}
`, getModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, data.EntityName)
}

func (g *ModuleGenerator) listUseCaseTemplate(data ModuleData) string {
	return fmt.Sprintf(`package usecases

import (
	"context"

	"%s/internal/modules/%s/domain"
)

type List%sUseCase struct {
	repo domain.%sRepository
}

func NewList%sUseCase(repo domain.%sRepository) *List%sUseCase {
	return &List%sUseCase{
		repo: repo,
	}
}

func (uc *List%sUseCase) Execute(ctx context.Context) ([]*domain.%s, error) {
	return uc.repo.List(ctx)
}
`, getModulePath(), data.ModuleName,
		data.ModuleNameTitle, data.EntityName,
		data.ModuleNameTitle, data.EntityName, data.ModuleNameTitle,
		data.ModuleNameTitle,
		data.ModuleNameTitle, data.EntityName)
}

func (g *ModuleGenerator) updateUseCaseTemplate(data ModuleData) string {
	return fmt.Sprintf(`package usecases

import (
	"context"
	"time"

	"%s/internal/modules/%s/application/dto"
	"%s/internal/modules/%s/domain"
)

type Update%sUseCase struct {
	repo domain.%sRepository
}

func NewUpdate%sUseCase(repo domain.%sRepository) *Update%sUseCase {
	return &Update%sUseCase{
		repo: repo,
	}
}

func (uc *Update%sUseCase) Execute(ctx context.Context, id int64, input dto.Update%sDTO) (*domain.%s, error) {
	%s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		%s.Name = input.Name
	}
	%s.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, %s); err != nil {
		return nil, err
	}

	return %s, nil
}
`, getModulePath(), data.ModuleName, getModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		strings.ToLower(data.EntityName[:1]),
		strings.ToLower(data.EntityName[:1]),
		strings.ToLower(data.EntityName[:1]),
		strings.ToLower(data.EntityName[:1]),
		strings.ToLower(data.EntityName[:1]))
}

func (g *ModuleGenerator) deleteUseCaseTemplate(data ModuleData) string {
	return fmt.Sprintf(`package usecases

import (
	"context"

	"%s/internal/modules/%s/domain"
)

type Delete%sUseCase struct {
	repo domain.%sRepository
}

func NewDelete%sUseCase(repo domain.%sRepository) *Delete%sUseCase {
	return &Delete%sUseCase{
		repo: repo,
	}
}

func (uc *Delete%sUseCase) Execute(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
`, getModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName)
}

func (g *ModuleGenerator) handlerTemplate(data ModuleData) string {
	return fmt.Sprintf(`package infra

import (
	"net/http"
	"strconv"

	"%s/internal/modules/%s/application/dto"
	"%s/internal/modules/%s/application/usecases"
	"github.com/gin-gonic/gin"
)

type %sHandler struct {
	createUC *usecases.Create%sUseCase
	getUC    *usecases.Get%sUseCase
	listUC   *usecases.List%sUseCase
	updateUC *usecases.Update%sUseCase
	deleteUC *usecases.Delete%sUseCase
}

func New%sHandler(
	createUC *usecases.Create%sUseCase,
	getUC *usecases.Get%sUseCase,
	listUC *usecases.List%sUseCase,
	updateUC *usecases.Update%sUseCase,
	deleteUC *usecases.Delete%sUseCase,
) *%sHandler {
	return &%sHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *%sHandler) Create(c *gin.Context) {
	var input dto.Create%sDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	%s, err := h.createUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, %s)
}

func (h *%sHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	%s, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, %s)
}

func (h *%sHandler) List(c *gin.Context) {
	%s, err := h.listUC.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, %s)
}

func (h *%sHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input dto.Update%sDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	%s, err := h.updateUC.Execute(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, %s)
}

func (h *%sHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.deleteUC.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
`, getModulePath(), data.ModuleName, getModulePath(), data.ModuleName,
		data.ModuleNameTitle,
		data.EntityName, data.EntityName, data.ModuleNameTitle,
		data.EntityName, data.EntityName,
		data.ModuleNameTitle,
		data.EntityName, data.EntityName, data.ModuleNameTitle,
		data.EntityName, data.EntityName,
		data.ModuleNameTitle, data.ModuleNameTitle,
		data.ModuleNameTitle, data.EntityName,
		strings.ToLower(data.ModuleName), strings.ToLower(data.ModuleName),
		data.ModuleNameTitle,
		strings.ToLower(data.ModuleName), strings.ToLower(data.ModuleName),
		data.ModuleNameTitle,
		data.ModuleName, data.ModuleName,
		data.ModuleNameTitle,
		data.EntityName,
		strings.ToLower(data.ModuleName), strings.ToLower(data.ModuleName),
		data.ModuleNameTitle)
}

func (g *ModuleGenerator) routesTemplate(data ModuleData) string {
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

func (g *ModuleGenerator) repositoryImplTemplate(data ModuleData) string {
	entityVar := strings.ToLower(data.EntityName[:1])
	return fmt.Sprintf(`package infra

import (
	"context"
	"database/sql"

	"%s/internal/modules/%s/domain"
)

type %sRepository struct {
	db *sql.DB
}

func New%sRepository(db *sql.DB) *%sRepository {
	return &%sRepository{
		db: db,
	}
}

func (r *%sRepository) Create(ctx context.Context, %s *domain.%s) error {
	query := `+"`"+`
		INSERT INTO %s (name, created_at, updated_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`+"`"+`

	err := r.db.QueryRowContext(
		ctx,
		query,
		%s.Name,
		%s.CreatedAt,
		%s.UpdatedAt,
	).Scan(&%s.ID)

	return err
}

func (r *%sRepository) GetByID(ctx context.Context, id int64) (*domain.%s, error) {
	query := `+"`"+`
		SELECT id, name, created_at, updated_at
		FROM %s
		WHERE id = $1
	`+"`"+`

	%s := &domain.%s{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&%s.ID,
		&%s.Name,
		&%s.CreatedAt,
		&%s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.Err%sNotFound
	}

	if err != nil {
		return nil, err
	}

	return %s, nil
}

func (r *%sRepository) List(ctx context.Context) ([]*domain.%s, error) {
	query := `+"`"+`
		SELECT id, name, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`+"`"+`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var %s []*domain.%s
	for rows.Next() {
		%s := &domain.%s{}
		if err := rows.Scan(
			&%s.ID,
			&%s.Name,
			&%s.CreatedAt,
			&%s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		%s = append(%s, %s)
	}

	return %s, rows.Err()
}

func (r *%sRepository) Update(ctx context.Context, %s *domain.%s) error {
	query := `+"`"+`
		UPDATE %s
		SET name = $1, updated_at = $2
		WHERE id = $3
	`+"`"+`

	_, err := r.db.ExecContext(
		ctx,
		query,
		%s.Name,
		%s.UpdatedAt,
		%s.ID,
	)

	return err
}

func (r *%sRepository) Delete(ctx context.Context, id int64) error {
	query := `+"`DELETE FROM %s WHERE id = $1`"+`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
`, getModulePath(), data.ModuleName,
		data.EntityName,
		data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, entityVar, data.EntityName,
		data.ModuleName,
		entityVar, entityVar, entityVar, entityVar,
		data.EntityName, data.EntityName,
		data.ModuleName,
		entityVar, data.EntityName,
		entityVar, entityVar, entityVar, entityVar,
		data.EntityName,
		entityVar,
		data.EntityName, data.EntityName,
		data.ModuleName,
		data.ModuleName, data.EntityName,
		entityVar, data.EntityName,
		entityVar, entityVar, entityVar, entityVar,
		data.ModuleName, data.ModuleName, entityVar,
		data.ModuleName,
		data.EntityName, entityVar, data.EntityName,
		data.ModuleName,
		entityVar, entityVar, entityVar,
		data.EntityName, data.ModuleName)
}

// Helper function to get module path from go.mod
func getModulePath() string {
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
