package templates

import (
	"fmt"
	"strings"
)

func ErrorsTemplate(data ModuleData) string {
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
