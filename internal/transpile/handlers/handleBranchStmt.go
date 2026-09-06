package handlers

import (
	"go/ast"
	"go/token"
)

// handleBranchStmt translates break and continue. A labelled one jumps
// somewhere C++ would have to name too, and goto has no place in a sketch this
// tool writes, so both stop the run.
func handleBranchStmt(stmt *ast.BranchStmt) string {
	if stmt.Label != nil {
		unsupported(stmt, "a labelled %s", stmt.Tok)
	}
	switch stmt.Tok {
	case token.BREAK:
		return "break;"
	case token.CONTINUE:
		return "continue;"
	}
	// FALLTHROUGH is handled where the case clause is written, and GOTO has no
	// translation.
	unsupported(stmt, "%s", stmt.Tok)
	return ""
}
