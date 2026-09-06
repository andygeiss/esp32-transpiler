package handlers

import (
	"go/ast"
	"strings"
)

func handleForStmt(stmt *ast.ForStmt) string {
	cond := ""
	if stmt.Cond != nil {
		cond = HandleExpr(stmt.Cond)
	}

	var code strings.Builder
	if stmt.Init == nil && stmt.Post == nil {
		// A Go for with only a condition is C++'s while, and one with nothing
		// at all loops forever.
		if cond == "" {
			cond = "true"
		}
		code.WriteString("while (" + cond + ") {")
	} else {
		code.WriteString("for (" + handleForHeader(stmt.Init) + ";" + cond + ";" + handleForHeader(stmt.Post) + ") {")
	}
	code.WriteString(handleBlockStmt(stmt.Body))
	code.WriteString("}")
	return code.String()
}

// handleForHeader writes the init or post clause of a for. They carry no
// terminator of their own: the header supplies the semicolons.
func handleForHeader(stmt ast.Stmt) string {
	switch s := stmt.(type) {
	case nil:
		return ""
	case *ast.AssignStmt:
		return handleAssignStmt(s)
	case *ast.IncDecStmt:
		return handleIncDecStmt(s)
	case *ast.ExprStmt:
		return handleExprStmt(s)
	}
	unsupported(stmt, "%s in a for header", describe(stmt))
	return ""
}
