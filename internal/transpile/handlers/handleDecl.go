package handlers

import "go/ast"

// HandleDecl turns one top-level declaration into its Arduino equivalent.
// Anything else transpiles to nothing.
func HandleDecl(decl ast.Decl) string {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return handleFuncDecl(d)
	case *ast.GenDecl:
		return handleGenDecl(d)
	}
	return ""
}
