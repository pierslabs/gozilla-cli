package templates

import (
	"fmt"
	"strings"
)

func RepositoryImplTemplate(data ModuleData) string {
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
`, GetModulePath(), data.ModuleName,
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
