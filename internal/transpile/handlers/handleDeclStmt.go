package handlers

import "go/ast"

func handleDeclStmt(stmt *ast.DeclStmt) string {
	gd, ok := stmt.Decl.(*ast.GenDecl)
	if !ok {
		unsupported(stmt.Decl, "%T inside a function", stmt.Decl)
	}
	return handleGenDecl(gd)
}
