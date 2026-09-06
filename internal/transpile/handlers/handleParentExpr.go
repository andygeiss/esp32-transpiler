package handlers

import "go/ast"

func handleParenExpr(expr *ast.ParenExpr) string {
	return "(" + HandleExpr(expr.X) + ")"
}
