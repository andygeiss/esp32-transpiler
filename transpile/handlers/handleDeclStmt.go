package handlers

import "go/ast"

func handleDeclStmt(stmt *ast.DeclStmt) string {
	code := ""
	switch decl := stmt.Decl.(type) {
	case *ast.GenDecl:
		code += handleGenDecl(decl)
	}
	return code
}
