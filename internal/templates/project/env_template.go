package templates

import "fmt"

func EnvExampleTemplate(data ProjectData) string {
	return fmt.Sprintf(`PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/%s?sslmode=disable
ENVIRONMENT=development
`, data.ProjectName)
}
