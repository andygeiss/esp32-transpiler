package handlers

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

// UnsupportedError names a Go construct outside the subset mapping.go covers,
// and where it is. The transpiler stops on one rather than write a sketch that
// would not compile, or — worse — would compile and do something else on the
// board than the Go said.
type UnsupportedError struct {
	Construct string // what was found, in Go's words
	Pos       token.Position
}

func (e *UnsupportedError) Error() string {
	return fmt.Sprintf("%s: %s is not supported", e.Pos, e.Construct)
}

// stop is what a handler raises when it meets something it cannot translate.
// It never leaves the package: HandleDecl recovers it and returns it as an
// UnsupportedError. The handlers are a recursive descent over the syntax tree
// and every one of them can meet an unsupported node, so threading an error
// out of each would bury the one line that does the translating.
type stop struct {
	construct string
	pos       token.Pos
}

// unsupported reports node as outside the subset and unwinds to HandleDecl.
func unsupported(node ast.Node, format string, a ...any) {
	pos := token.NoPos
	if node != nil {
		pos = node.Pos()
	}
	panic(&stop{construct: fmt.Sprintf(format, a...), pos: pos})
}

// describe names a node the way a Go programmer would, so the message says
// what to take out of the source rather than which syntax type it parsed as.
func describe(n ast.Node) string {
	switch x := n.(type) {
	case *ast.FuncLit:
		return "a function literal"
	case *ast.DeferStmt:
		return "a defer statement"
	case *ast.GoStmt:
		return "a go statement"
	case *ast.LabeledStmt:
		return "a label"
	case *ast.RangeStmt:
		return "a range loop"
	case *ast.SelectStmt:
		return "a select statement"
	case *ast.SendStmt:
		return "a channel send"
	case *ast.TypeSwitchStmt:
		return "a type switch"
	case *ast.BlockStmt:
		return "a bare block"
	case ast.Expr: // last: every case above it would be unreachable after it
		return types.ExprString(x)
	}
	return fmt.Sprintf("%T", n)
}
