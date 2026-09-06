package handlers

import "go/ast"

func HandleExpr(expr ast.Expr) string {
	code := ""
	switch e := expr.(type) {
	case *ast.BasicLit:
		code += handleBasicLit(e)
	case *ast.BinaryExpr:
		code += handleBinaryExpr(e)
	case *ast.CallExpr:
		code += handleCallExpr(e)
	case *ast.Ident:
		code += handleIdent(e)
	case *ast.ParenExpr:
		code += handleParenExpr(e)
	case *ast.SelectorExpr:
		code += handleSelectorExpr(e)
	}
	return code
}
