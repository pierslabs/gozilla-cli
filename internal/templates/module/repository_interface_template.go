package templates

import (
	"fmt"
	"strings"
)

func RepositoryInterfaceTemplate(data ModuleData) string {
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
