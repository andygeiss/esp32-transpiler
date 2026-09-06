package handlers

import "go/ast"

func handleBinaryExpr(expr ast.Expr) string {
	be := expr.(*ast.BinaryExpr)
	code := HandleExpr(be.X)
	code += be.Op.String()
	code += HandleExpr(be.Y)
	return code
}
