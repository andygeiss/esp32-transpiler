package handlers

import (
	"go/ast"
	"strings"
)

func handleFuncDeclParams(t *ast.FuncType) string {
	if t.Params == nil || len(t.Params.List) == 0 {
		return ""
	}
	values := make([]string, 0, len(t.Params.List))
	for _, field := range t.Params.List {
		ftype := handleType(field.Type)
		if len(field.Names) == 0 {
			// C++ lets a parameter go unnamed too.
			values = append(values, ftype)
			continue
		}
		for _, name := range field.Names {
			values = append(values, ftype+" "+name.Name)
		}
	}
	return strings.Join(values, ",")
}
