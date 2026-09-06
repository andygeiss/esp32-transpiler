package handlers

import (
	"fmt"
	"go/ast"
)

func handleIfStmt(stmt *ast.IfStmt) string {
	cond := HandleExpr(stmt.Cond)
	body := handleBlockStmt(stmt.Body)
	code := fmt.Sprintf(`if (%s) { %s }`, cond, body)
	if stmt.Else != nil {
		tail := handleBlockStmt(stmt.Else.(*ast.BlockStmt))
		code += fmt.Sprintf(" else { %s }", tail)
	}
	return code
}
