package handlers

import "go/ast"

func handleParenExpr(stmt *ast.ParenExpr) string {
	code := ""
	code += HandleExpr(stmt.X)
	return code
}
