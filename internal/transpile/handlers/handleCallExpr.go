package handlers

import (
	"go/ast"
	"go/token"
	"strings"
)

func handleCallExpr(expr *ast.CallExpr) string {
	if expr.Ellipsis != token.NoPos {
		// Dropping the ... would pass the collection where the callee wants
		// its elements.
		unsupported(expr, "a call spreading its last argument with ...")
	}
	args := make([]string, 0, len(expr.Args))
	for _, arg := range expr.Args {
		args = append(args, HandleExpr(arg))
	}
	return HandleExpr(expr.Fun) + "(" + strings.Join(args, ",") + ")"
}
