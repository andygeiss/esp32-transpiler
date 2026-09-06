package handlers

import (
	"go/ast"
	"go/token"
)

func handleGenDecl(decl ast.Decl) string {
	gd := decl.(*ast.GenDecl)
	code := ""
	switch gd.Tok {
	case token.CONST:
		code += "const "
	case token.VAR:
		code += ""
	}
	code += handleSpecs(gd.Specs)
	return code
}
