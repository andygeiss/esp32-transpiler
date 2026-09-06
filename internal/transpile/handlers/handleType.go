package handlers

import "go/ast"

// handleType writes a type: the plain names Go and C++ share, and the package
// ones mapping.go knows. Slices, maps, pointers and channels have no place in
// a sketch this tool writes, so they stop the run.
func handleType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return handleIdent(t)
	case *ast.SelectorExpr:
		return handleSelectorExpr(t)
	}
	unsupported(expr, "the type %s", describe(expr))
	return ""
}
