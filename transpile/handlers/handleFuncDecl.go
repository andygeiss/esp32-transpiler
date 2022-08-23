package handlers

import "go/ast"

func handleFuncDecl(decl ast.Decl) string {
	fd := decl.(*ast.FuncDecl)
	code := ""
	name := ""
	code += handleFuncDeclType(fd.Type)
	code += " "
	name = handleFuncDeclName(fd.Name)
	if name == "NewController" {
		return ""
	}
	code += name
	code += "("
	code += handleFuncDeclParams(fd.Type)
	code += ") {"
	code += handleBlockStmt(fd.Body)
	code += "}"
	return code
}
