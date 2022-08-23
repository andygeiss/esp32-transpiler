package handlers

import "go/ast"

func handleValueSpecValues(values []ast.Expr) string {
	code := ""
	for _, value := range values {
		switch v := value.(type) {
		case *ast.BasicLit:
			code += handleBasicLit(v)
		}
	}
	return code
}
