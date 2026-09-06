package handlers

import (
	"go/ast"
	"strings"
)

// handleSpecs writes the specs of one declaration. prefix goes on each of
// them, because "const ( a = 1; b = 2 )" declares two constants and C++ wants
// the keyword on both.
func handleSpecs(specs []ast.Spec, prefix string) string {
	var code strings.Builder
	for _, spec := range specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			code.WriteString(handleImportSpec(s))
		case *ast.ValueSpec:
			code.WriteString(handleValueSpec(s, prefix) + ";")
		default:
			unsupported(spec, "%s", describe(spec))
		}
	}
	return code.String()
}
