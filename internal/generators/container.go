package generators

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	templates "github.com/pierslabs/gozilla/internal/templates/module"
)

type ContainerUpdater struct{}

func NewContainerUpdater() *ContainerUpdater {
	return &ContainerUpdater{}
}

func (u *ContainerUpdater) AddModule(moduleName string) error {
	containerPath := filepath.Join("internal", "infrastructure", "container", "container.go")

	// Read the file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, containerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse container.go: %w", err)
	}

	moduleNameTitle := strings.Title(moduleName)
	moduleVarName := moduleNameTitle + "Module"
	moduleImportPath := fmt.Sprintf("%s/internal/modules/%s", templates.GetModulePath(), moduleName)

	// Add import
	u.addImport(file, moduleName, moduleImportPath)

	// Update Container struct
	u.addFieldToStruct(file, "Container", moduleVarName, "*"+moduleName+"."+moduleNameTitle+"Module")

	// Update NewContainer function
	u.addModuleToConstructor(file, moduleVarName, moduleName, moduleNameTitle)

	// Update RegisterRoutes method
	u.addRouteRegistration(file, moduleVarName)

	// Write back to file
	f, err := os.Create(containerPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if err := format.Node(f, fset, file); err != nil {
		return fmt.Errorf("failed to format file: %w", err)
	}

	return nil
}

func (u *ContainerUpdater) addImport(file *ast.File, alias, path string) {
	// Check if import already exists
	for _, imp := range file.Imports {
		if imp.Path.Value == `"`+path+`"` {
			return
		}
	}

	// Add new import
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + path + `"`,
		},
	}

	// Find or create import declaration
	var importDecl *ast.GenDecl
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			importDecl = genDecl
			break
		}
	}

	if importDecl == nil {
		// Create new import declaration
		importDecl = &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: []ast.Spec{},
		}
		file.Decls = append([]ast.Decl{importDecl}, file.Decls...)
	}

	importDecl.Specs = append(importDecl.Specs, newImport)
}

func (u *ContainerUpdater) addFieldToStruct(file *ast.File, structName, fieldName, fieldType string) {
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// Check if field already exists
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					if name.Name == fieldName {
						return // Field already exists
					}
				}
			}

			// Add new field
			newField := &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fieldName)},
				Type:  parseTypeExpr(fieldType),
			}

			structType.Fields.List = append(structType.Fields.List, newField)
			return
		}
	}
}

func (u *ContainerUpdater) addModuleToConstructor(file *ast.File, moduleVarName, moduleName, moduleNameTitle string) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "NewContainer" {
			continue
		}

		// Find the return statement
		for _, stmt := range funcDecl.Body.List {
			retStmt, ok := stmt.(*ast.ReturnStmt)
			if !ok {
				continue
			}

			for _, result := range retStmt.Results {
				composite, ok := result.(*ast.UnaryExpr)
				if !ok {
					continue
				}

				compositeLit, ok := composite.X.(*ast.CompositeLit)
				if !ok {
					continue
				}

				// Check if field already exists
				for _, elt := range compositeLit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == moduleVarName {
						return // Already exists
					}
				}

				// Add new field initialization
				newElt := &ast.KeyValueExpr{
					Key: ast.NewIdent(moduleVarName),
					Value: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(moduleName),
							Sel: ast.NewIdent("New" + moduleNameTitle + "Module"),
						},
						Args: []ast.Expr{ast.NewIdent("db")},
					},
				}

				compositeLit.Elts = append(compositeLit.Elts, newElt)
				return
			}
		}
	}
}

func (u *ContainerUpdater) addRouteRegistration(file *ast.File, moduleVarName string) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != "RegisterRoutes" {
			continue
		}

		// Check if registration already exists
		for _, stmt := range funcDecl.Body.List {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
					if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
						if x, ok := selExpr.X.(*ast.SelectorExpr); ok {
							if x.Sel.Name == moduleVarName {
								return // Already registered
							}
						}
					}
				}
			}
		}

		// Add route registration call
		newStmt := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("c"),
						Sel: ast.NewIdent(moduleVarName),
					},
					Sel: ast.NewIdent("RegisterRoutes"),
				},
				Args: []ast.Expr{ast.NewIdent("api")},
			},
		}

		funcDecl.Body.List = append(funcDecl.Body.List, newStmt)
		return
	}
}

func parseTypeExpr(typeStr string) ast.Expr {
	// Simple type expression parser
	if strings.HasPrefix(typeStr, "*") {
		return &ast.StarExpr{
			X: parseTypeExpr(typeStr[1:]),
		}
	}

	if strings.Contains(typeStr, ".") {
		parts := strings.Split(typeStr, ".")
		return &ast.SelectorExpr{
			X:   ast.NewIdent(parts[0]),
			Sel: ast.NewIdent(parts[1]),
		}
	}

	return ast.NewIdent(typeStr)
}
