package handlers

import "go/ast"

func handleValueSpecType(expr ast.Expr) string {
	code := ""
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		code += handleSelectorExpr(t)
	case *ast.Ident:
		code += handleIdent(t)
	}
	return code
}
