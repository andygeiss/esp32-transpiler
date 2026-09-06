package handlers

import "go/ast"

// handleReturnStmt carries a return across. "return nil" is how a Go
// controller says nothing went wrong, and the void function it becomes has
// nothing to give back — but it still stops there, so it stays a return.
// handleFuncDecl drops the one that ends the function, where it says nothing
// the closing brace does not.
func handleReturnStmt(stmt *ast.ReturnStmt) string {
	switch len(stmt.Results) {
	case 0:
		return "return;"
	case 1:
		if isNothing(stmt.Results[0]) {
			return "return;"
		}
		return "return " + HandleExpr(stmt.Results[0]) + ";"
	}
	unsupported(stmt, "a return of more than one value")
	return ""
}

// endsFunction reports whether the return says no more than the closing brace
// of a void function already does.
func endsFunction(stmt *ast.ReturnStmt) bool {
	switch len(stmt.Results) {
	case 0:
		return true
	case 1:
		return isNothing(stmt.Results[0])
	}
	return false
}

func isNothing(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}
