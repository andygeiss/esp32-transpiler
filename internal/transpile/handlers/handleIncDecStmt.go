package handlers

import "go/ast"

// handleIncDecStmt writes x++ or x-- without its terminator, so a for header
// can reuse it.
func handleIncDecStmt(stmt *ast.IncDecStmt) string {
	return HandleExpr(stmt.X) + stmt.Tok.String()
}
