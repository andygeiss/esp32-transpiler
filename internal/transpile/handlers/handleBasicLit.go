package handlers

import "go/ast"

func handleBasicLit(bl *ast.BasicLit) string {
	return bl.Value
}
