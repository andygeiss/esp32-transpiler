package handlers

import "go/ast"

func handleSwitchStmt(stmt *ast.SwitchStmt) string {
	if stmt.Init != nil {
		unsupported(stmt, "a switch with an init statement")
	}
	// C++ switches on a value; a Go switch without one is a chain of ifs.
	if stmt.Tag == nil {
		unsupported(stmt, "a switch without a value")
	}
	code := "switch ("
	code += HandleExpr(stmt.Tag)
	code += "){"
	code += handleBlockStmt(stmt.Body)
	code += "}"
	return code
}
