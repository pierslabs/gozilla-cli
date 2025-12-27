package templates

import "fmt"

func CreateDTOTemplate(data ModuleData) string {
	return fmt.Sprintf(`package dto

type Create%sDTO struct {
	Name string `+"`json:\"name\" binding:\"required\"`"+`
}
`, data.EntityName)
}

func UpdateDTOTemplate(data ModuleData) string {
	return fmt.Sprintf(`package dto

type Update%sDTO struct {
	Name string `+"`json:\"name\"`"+`
}
`, data.EntityName)
}
