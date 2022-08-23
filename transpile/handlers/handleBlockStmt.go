package handlers

import "go/ast"

func handleBlockStmt(body *ast.BlockStmt) string {
	code := ""
	if body == nil {
		return code
	}
	for _, stmt := range body.List {
		code += handleStmt(stmt)
	}
	return code
}
