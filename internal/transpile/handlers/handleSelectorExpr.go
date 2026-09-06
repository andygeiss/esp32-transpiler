package handlers

import "go/ast"

func handleSelectorExpr(expr ast.Expr) string {
	s := expr.(*ast.SelectorExpr)
	code := ""
	switch x := s.X.(type) {
	case *ast.Ident:
		code += handleIdent(x)
	}
	code += "."
	code += handleIdent(s.Sel)
	if val, ok := mapping[code]; ok {
		code = val
	}
	return code
}
