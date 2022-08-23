package handlers

import "go/ast"

func handleStmt(stmt ast.Stmt) string {
	code := ""
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		code += handleAssignStmt(s)
		code += ";"
	case *ast.BranchStmt:
		code += handleBranchStmt()
	case *ast.CaseClause:
		code += handleCaseClause(s)
	case *ast.DeclStmt:
		code += handleDeclStmt(s)
	case *ast.ExprStmt:
		code += handleExprStmt(s)
		code += ";"
	case *ast.ForStmt:
		code += handleForStmt(s)
	case *ast.IfStmt:
		code += handleIfStmt(s)
	case *ast.SwitchStmt:
		code += handleSwitchStmt(s)
	}
	return code
}
