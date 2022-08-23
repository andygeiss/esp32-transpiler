package handlers

import (
	"go/ast"
	"strings"
)

func handleCaseClause(cc *ast.CaseClause) string {
	code := "case "
	clauses := make([]string, 0)
	for _, clause := range cc.List {
		clauses = append(clauses, HandleExpr(clause))
	}
	code += strings.Join(clauses, ",")
	code += ":"
	for _, body := range cc.Body {
		code += handleStmt(body)
	}
	return code
}
