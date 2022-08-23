package handlers

import "go/ast"

func handleForStmt(stmt *ast.ForStmt) string {
	code := ""
	if stmt.Init == nil && stmt.Post == nil {
		code += "while"
	} else {
		code += "for"
	}
	code += "("                         // stmt.Init
	code += handleBinaryExpr(stmt.Cond) // stmt.Cond
	code += ""                          // stmt.Post
	code += ") {"
	code += handleBlockStmt(stmt.Body) // stmt.Body
	code += "}"
	return code
}
