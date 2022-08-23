package handlers

import "go/ast"

func handleSpecs(specs []ast.Spec) string {
	code := ""
	for _, spec := range specs {
		switch spec.(type) {
		case *ast.ImportSpec:
			code += handleImportSpec(spec)
		case *ast.ValueSpec:
			code += handleValueSpec(spec) + ";"
		}
	}
	return code
}
