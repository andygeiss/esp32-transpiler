package handlers

import (
	"go/ast"
)

func handleAssignStmt(as *ast.AssignStmt) string {
	code := handleAssignStmtExpr(as.Lhs)
	code += as.Tok.String()
	code += handleAssignStmtExpr(as.Rhs)
	return code
}
