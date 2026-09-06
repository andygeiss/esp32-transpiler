package handlers

import (
	"go/ast"
	"go/token"
)

func handleGenDecl(gd *ast.GenDecl) string {
	prefix := ""
	switch gd.Tok {
	case token.CONST:
		prefix = "const "
	case token.IMPORT, token.VAR:
	case token.TYPE:
		unsupported(gd, "a type declaration")
	}
	return handleSpecs(gd.Specs, prefix)
}
