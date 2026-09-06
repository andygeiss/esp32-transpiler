package handlers

import "go/ast"

func handleIdent(ident *ast.Ident) string {
	switch ident.Name {
	case "string":
		// const, because what a Go string holds is a literal in the sketch's
		// flash and C++ will not hand one out as a writable char*.
		return "const char*"
	case "nil":
		return "nullptr"
	}
	return ident.Name
}
