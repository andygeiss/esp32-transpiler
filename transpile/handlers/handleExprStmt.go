package handlers

import "go/ast"

func handleExprStmt(stmt *ast.ExprStmt) string {
	code := ""
	switch x := stmt.X.(type) {
	case *ast.CallExpr:
		code += handleCallExpr(x)
	}
	return code
}
