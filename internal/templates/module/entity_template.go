package templates

import "fmt"

func EntityTemplate(data ModuleData) string {
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
