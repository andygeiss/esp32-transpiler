package handlers

import "go/ast"

// handleFuncDeclType writes what a function gives back. Setup and Loop never
// reach here: their mapped names carry the void already.
func handleFuncDeclType(t *ast.FuncType) string {
	if t.Results == nil || len(t.Results.List) == 0 {
		return "void"
	}
	results := 0
	for _, field := range t.Results.List {
		results += max(len(field.Names), 1)
	}
	if results > 1 {
		unsupported(t, "a function returning more than one value")
	}
	result := t.Results.List[0].Type
	if ident, ok := result.(*ast.Ident); ok && ident.Name == "error" {
		// Only Setup and Loop may return one, and they map to void.
		unsupported(t, "an error result on a function other than Setup or Loop")
	}
	return handleType(result)
}
