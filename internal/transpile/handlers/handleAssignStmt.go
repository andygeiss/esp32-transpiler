package handlers

import (
	"go/ast"
	"go/token"
)

// handleAssignStmt writes an assignment without its terminator, so a for
// header can reuse it.
func handleAssignStmt(as *ast.AssignStmt) string {
	// C++ assigns one thing at a time; "x, y = 1, 2" has no counterpart.
	if len(as.Lhs) > 1 || len(as.Rhs) > 1 {
		unsupported(as, "an assignment to more than one variable")
	}
	// "_ = f()" runs f and throws the answer away, which C++ does by writing
	// the call on its own.
	if ident, ok := as.Lhs[0].(*ast.Ident); ok && ident.Name == "_" && as.Tok != token.DEFINE {
		return HandleExpr(as.Rhs[0])
	}
	code := HandleExpr(as.Lhs[0])
	if as.Tok == token.AND_NOT_ASSIGN {
		// Go's x &^= y clears bits; C++ spells that x &= ~(y).
		return code + "&=~(" + HandleExpr(as.Rhs[0]) + ")"
	}
	if as.Tok == token.DEFINE {
		// ":=" declares, and auto lets the sketch compiler work the type out
		// the way Go does.
		return "auto " + code + " = " + HandleExpr(as.Rhs[0])
	}
	code += as.Tok.String()
	code += HandleExpr(as.Rhs[0])
	return code
}
