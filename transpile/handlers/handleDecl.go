package handlers

import "go/ast"

func HandleDecl(decl ast.Decl, dst chan<- string, done chan<- bool) {
	code := ""
	switch d := decl.(type) {
	case *ast.FuncDecl:
		code += handleFuncDecl(d)
	case *ast.GenDecl:
		code += handleGenDecl(d)
	}
	dst <- code
	close(dst)
	done <- true
}
