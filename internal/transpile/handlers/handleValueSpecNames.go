package handlers

import (
	"go/ast"
	"strings"
)

func handleValueSpecNames(names []*ast.Ident) string {
	var code strings.Builder
	for _, name := range names {
		code.WriteString(handleIdent(name))
	}
	return code.String()
}
