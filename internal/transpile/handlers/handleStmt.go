package handlers

import "go/ast"

func handleStmt(stmt ast.Stmt) string {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return handleAssignStmt(s) + ";"
	case *ast.BranchStmt:
		return handleBranchStmt(s)
	case *ast.CaseClause:
		return handleCaseClause(s)
	case *ast.DeclStmt:
		return handleDeclStmt(s)
	case *ast.ExprStmt:
		return handleExprStmt(s) + ";"
	case *ast.ForStmt:
		return handleForStmt(s)
	case *ast.IfStmt:
		return handleIfStmt(s)
	case *ast.IncDecStmt:
		return handleIncDecStmt(s) + ";"
	case *ast.ReturnStmt:
		return handleReturnStmt(s)
	case *ast.SwitchStmt:
		return handleSwitchStmt(s)
	case *ast.EmptyStmt:
		return ""
	}
	unsupported(stmt, "%s", describe(stmt))
	return ""
}
