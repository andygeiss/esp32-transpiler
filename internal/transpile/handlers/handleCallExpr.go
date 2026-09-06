package handlers

import (
	"go/ast"
	"strings"
)

func handleCallExpr(expr *ast.CallExpr) string {
	code := HandleExpr(expr.Fun)
	code += "("
	args := make([]string, 0)
	for _, arg := range expr.Args {
		args = append(args, HandleExpr(arg))
	}
	code += strings.Join(args, ",")
	code += ")"
	return code
}
