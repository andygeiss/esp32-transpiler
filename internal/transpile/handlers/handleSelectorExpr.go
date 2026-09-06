package handlers

import "go/ast"

// handleSelectorExpr writes a.B, whatever a is, and then asks mapping.go
// whether the whole of it is a name the Arduino spells differently.
func handleSelectorExpr(s *ast.SelectorExpr) string {
	code := HandleExpr(s.X) + "." + handleIdent(s.Sel)
	if val, ok := mapping[code]; ok {
		return val
	}
	return code
}
