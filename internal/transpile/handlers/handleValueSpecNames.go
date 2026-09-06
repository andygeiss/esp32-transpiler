package handlers

import "go/ast"

func handleValueSpecNames(names []*ast.Ident) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		out = append(out, name.Name)
	}
	return out
}
