package handlers

import (
	"go/ast"
	"go/token"
)

// HandleDecl turns one top-level declaration into its Arduino equivalent. It
// is the boundary the handlers report through: anything they cannot translate
// comes back as an *UnsupportedError naming the construct and its line, and
// nothing is written for that declaration.
func HandleDecl(fset *token.FileSet, decl ast.Decl) (code string, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		s, ok := r.(*stop)
		if !ok {
			panic(r) // a real bug: let it crash
		}
		code, err = "", &UnsupportedError{Construct: s.construct, Pos: fset.Position(s.pos)}
	}()

	switch d := decl.(type) {
	case *ast.FuncDecl:
		return handleFuncDecl(d), nil
	case *ast.GenDecl:
		return handleGenDecl(d), nil
	}
	unsupported(decl, "%T", decl)
	return "", nil // unreachable: unsupported panics
}
