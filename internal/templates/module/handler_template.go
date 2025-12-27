package templates

import (
	"fmt"
	"strings"
)

func HandlerTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName, GetModulePath(), data.ModuleName,
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
