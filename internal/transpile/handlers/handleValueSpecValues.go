package handlers

import (
	"go/ast"
	"strings"
)

func handleValueSpecValues(values []ast.Expr) string {
	var code strings.Builder
	for _, value := range values {
		switch v := value.(type) {
		case *ast.BasicLit:
			code.WriteString(handleBasicLit(v))
		}
	}
	return code.String()
}
