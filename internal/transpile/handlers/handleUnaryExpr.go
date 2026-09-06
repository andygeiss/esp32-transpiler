package handlers

import (
	"go/ast"
	"go/token"
)

func handleUnaryExpr(expr *ast.UnaryExpr) string {
	op := ""
	switch expr.Op {
	case token.NOT, token.SUB, token.ADD:
		op = expr.Op.String()
	case token.XOR:
		op = "~" // Go spells the bitwise complement ^x, C++ ~x
	default:
		// & takes an address and <- reads a channel; neither is in the subset.
		unsupported(expr, "the unary operator %s", expr.Op)
	}
	return op + HandleExpr(expr.X)
}
