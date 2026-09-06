package handlers

import (
	"go/ast"
	"strings"
)

func handleAssignStmtExpr(e []ast.Expr) string {
	ops := make([]string, 0)
	code := ""
	for _, op := range e {
		ops = append(ops, HandleExpr(op))
	}
	code += strings.Join(ops, ",")
	return code
}
