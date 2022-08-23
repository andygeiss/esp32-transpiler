package handlers

import "go/ast"

func handleValueSpecNames(names []*ast.Ident) string {
	code := ""
	for _, name := range names {
		code += handleIdent(name)
	}
	return code
}
