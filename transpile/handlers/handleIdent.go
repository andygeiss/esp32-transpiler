package handlers

import "go/ast"

func handleIdent(expr ast.Expr) string {
	ident := expr.(*ast.Ident)
	code := ""
	switch ident.Name {
	case "string":
		code += "char*"
	default:
		code += ident.Name
	}
	return code
}
