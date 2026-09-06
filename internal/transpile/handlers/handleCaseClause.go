package handlers

import (
	"go/ast"
	"go/token"
	"strings"
)

// handleCaseClause writes one case of a switch. Two things differ from Go and
// both change what the board does, so both are translated rather than copied:
// a C++ case falls into the next one unless it breaks, where a Go case never
// does, and C++ has no "case 1, 2:" — it wants a label each.
func handleCaseClause(cc *ast.CaseClause) string {
	var code strings.Builder
	if len(cc.List) == 0 {
		code.WriteString("default:")
	} else {
		for _, clause := range cc.List {
			code.WriteString("case ")
			code.WriteString(HandleExpr(clause))
			code.WriteString(":")
		}
	}

	body := cc.Body
	// Go's fallthrough is what C++ does on its own, so it becomes the absence
	// of the break below.
	falls := endsInFallthrough(body)
	if falls {
		body = body[:len(body)-1]
	}
	for _, stmt := range body {
		code.WriteString(handleStmt(stmt))
	}
	if !falls && !endsInJump(body) {
		code.WriteString("break;")
	}
	return code.String()
}

func endsInFallthrough(body []ast.Stmt) bool {
	if len(body) == 0 {
		return false
	}
	b, ok := body[len(body)-1].(*ast.BranchStmt)
	return ok && b.Tok == token.FALLTHROUGH
}

// endsInJump reports whether the case already leaves the switch on its own, so
// the break would be dead code.
func endsInJump(body []ast.Stmt) bool {
	if len(body) == 0 {
		return false
	}
	switch s := body[len(body)-1].(type) {
	case *ast.BranchStmt:
		return s.Tok == token.BREAK || s.Tok == token.CONTINUE
	case *ast.ReturnStmt:
		return true
	}
	return false
}
