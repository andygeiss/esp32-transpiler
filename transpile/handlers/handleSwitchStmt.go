package handlers

import "go/ast"

func handleSwitchStmt(stmt *ast.SwitchStmt) string {
	code := "switch ("
	code += HandleExpr(stmt.Tag)
	code += "){"
	code += handleBlockStmt(stmt.Body)
	code += "}"
	return code
}
