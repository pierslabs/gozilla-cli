package templates

import (
	"fmt"
	"strings"
)

func CreateUseCaseTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName, GetModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		entityVar, data.EntityName,
		entityVar,
		entityVar)
}

func GetUseCaseTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName, data.EntityName)
}

func ListUseCaseTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName,
		data.ModuleNameTitle, data.EntityName,
		data.ModuleNameTitle, data.EntityName, data.ModuleNameTitle,
		data.ModuleNameTitle,
		data.ModuleNameTitle, data.EntityName)
}

func UpdateUseCaseTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName, GetModulePath(), data.ModuleName,
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

func DeleteUseCaseTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName,
		data.EntityName, data.EntityName,
		data.EntityName, data.EntityName, data.EntityName,
		data.EntityName,
		data.EntityName)
}
