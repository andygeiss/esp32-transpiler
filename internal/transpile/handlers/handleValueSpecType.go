package handlers

import (
	"go/ast"
	"go/token"
)

// handleValueSpecType writes the type of a declaration. Go lets it go unsaid
// and C++ does too, but naming the one a literal implies keeps the sketch
// readable.
func handleValueSpecType(s *ast.ValueSpec) string {
	if s.Type != nil {
		return handleType(s.Type)
	}
	if len(s.Values) == 0 {
		unsupported(s, "a declaration with neither a type nor a value")
	}
	if lit, ok := s.Values[0].(*ast.BasicLit); ok && len(s.Names) == 1 {
		return literalType(lit)
	}
	return "auto"
}

func literalType(lit *ast.BasicLit) string {
	switch lit.Kind {
	case token.INT:
		return "int"
	case token.FLOAT:
		return "float"
	case token.CHAR:
		return "char"
	case token.STRING:
		return "const char*"
	}
	unsupported(lit, "the literal %s", lit.Value)
	return ""
}
