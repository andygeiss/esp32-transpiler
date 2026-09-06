package handlers

import (
	"go/ast"
	"strings"
)

func handleCaseClause(cc *ast.CaseClause) string {
	var code strings.Builder
	code.WriteString("case ")
	clauses := make([]string, 0)
	for _, clause := range cc.List {
		clauses = append(clauses, HandleExpr(clause))
	}
	code.WriteString(strings.Join(clauses, ","))
	code.WriteString(":")
	for _, body := range cc.Body {
		code.WriteString(handleStmt(body))
	}
	return code.String()
}
