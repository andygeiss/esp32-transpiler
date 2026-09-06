package handlers

import "go/ast"

// handleExprStmt writes a call without its terminator, so a for header can
// reuse it. A statement that is any other expression has no effect in Go and
// the parser is the only thing that lets it through.
func handleExprStmt(stmt *ast.ExprStmt) string {
	call, ok := stmt.X.(*ast.CallExpr)
	if !ok {
		unsupported(stmt, "%s as a statement", describe(stmt.X))
	}
	return handleCallExpr(call)
}
