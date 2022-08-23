package handlers

import (
	"go/ast"
	"strings"
)

func handleFuncDeclParams(t *ast.FuncType) string {
	code := ""
	if t.Params == nil || t.Params.List == nil {
		return code
	}
	values := make([]string, 0)
	for _, field := range t.Params.List {
		ftype := ""
		switch ft := field.Type.(type) {
		case *ast.Ident:
			ftype = handleIdent(ft)
		}
		for _, names := range field.Names {
			values = append(values, ftype+" "+names.Name)
		}
	}
	code += strings.Join(values, ",")
	return code
}
