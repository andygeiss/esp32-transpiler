package handlers

import "go/ast"

func HandleExpr(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return handleBasicLit(e)
	case *ast.BinaryExpr:
		return handleBinaryExpr(e)
	case *ast.CallExpr:
		return handleCallExpr(e)
	case *ast.Ident:
		return handleIdent(e)
	case *ast.ParenExpr:
		return handleParenExpr(e)
	case *ast.SelectorExpr:
		return handleSelectorExpr(e)
	case *ast.UnaryExpr:
		return handleUnaryExpr(e)
	}
	unsupported(expr, "%s", describe(expr))
	return ""
}
