package handlers

import (
	"fmt"
	"go/ast"
)

func handleIfStmt(stmt *ast.IfStmt) string {
	if stmt.Init != nil {
		// The init would have to be lifted out of the if, and that changes
		// what the name it declares is in scope for.
		unsupported(stmt, "an if with an init statement")
	}
	code := fmt.Sprintf("if (%s) { %s }", HandleExpr(stmt.Cond), handleBlockStmt(stmt.Body))
	switch e := stmt.Else.(type) {
	case nil:
	case *ast.BlockStmt:
		code += fmt.Sprintf(" else { %s }", handleBlockStmt(e))
	case *ast.IfStmt:
		code += " else " + handleIfStmt(e)
	default:
		code += fmt.Sprintf(" else { %s }", handleStmt(e))
	}
	return code
}
