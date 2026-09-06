package handlers

import (
	"go/ast"
	"strings"
)

// handleValueSpec writes one declaration. prefix is the const a Go const asks
// for, left off when the type carries one already: C++ has no "const const".
func handleValueSpec(s *ast.ValueSpec, prefix string) string {
	rejectIota(s)
	typ := handleValueSpecType(s)
	names := handleValueSpecNames(s.Names)
	if !strings.HasPrefix(typ, "const") {
		typ = prefix + typ
	}

	if len(s.Values) == 0 {
		return typ + " " + strings.Join(names, ",")
	}
	// One value per name, or C++ has no way to spread them.
	if len(s.Values) != len(s.Names) {
		unsupported(s, "a declaration spreading %d values over %d names", len(s.Values), len(s.Names))
	}
	pairs := make([]string, 0, len(names))
	for i, name := range names {
		pairs = append(pairs, name+" = "+HandleExpr(s.Values[i]))
	}
	return typ + " " + strings.Join(pairs, ",")
}

// rejectIota stops a const that counts. C++ has no iota, and guessing what
// each line would come to is not something this tool does.
func rejectIota(s *ast.ValueSpec) {
	for _, value := range s.Values {
		ast.Inspect(value, func(n ast.Node) bool {
			if ident, ok := n.(*ast.Ident); ok && ident.Name == "iota" {
				unsupported(ident, "iota")
			}
			return true
		})
	}
}
