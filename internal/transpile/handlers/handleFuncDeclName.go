package handlers

import "go/ast"

func handleFuncDeclName(ident *ast.Ident) string {
	code := ""
	if ident == nil {
		return code
	}
	code += ident.Name
	if val, ok := mapping[code]; ok {
		code = val
	}
	return code
}
