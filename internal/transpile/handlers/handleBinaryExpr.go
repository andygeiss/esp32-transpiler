package handlers

import (
	"go/ast"
	"go/token"
)

func handleBinaryExpr(be *ast.BinaryExpr) string {
	switch {
	case be.Op == token.AND_NOT:
		// Go's a &^ b clears bits; C++ spells that a & ~b.
		return HandleExpr(be.X) + "&~" + HandleExpr(be.Y)
	case be.Op == token.ADD && (isStringLit(be.X) || isStringLit(be.Y)):
		// Go's + joins strings. C++ would add to a pointer instead, quietly,
		// and the sketch would print whatever that address held.
		unsupported(be, "joining strings with +")
	}
	return HandleExpr(be.X) + be.Op.String() + HandleExpr(be.Y)
}

func isStringLit(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}
