package handlers

import "go/ast"

func handleFuncDeclType(t *ast.FuncType) string {
	code := ""
	if t.Results == nil {
		code = "void"
	}
	return code
}
