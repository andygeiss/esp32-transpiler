package handlers

import (
	"go/ast"
	"strings"
)

func handleSpecs(specs []ast.Spec) string {
	var code strings.Builder
	for _, spec := range specs {
		switch spec.(type) {
		case *ast.ImportSpec:
			code.WriteString(handleImportSpec(spec))
		case *ast.ValueSpec:
			code.WriteString(handleValueSpec(spec) + ";")
		}
	}
	return code.String()
}
