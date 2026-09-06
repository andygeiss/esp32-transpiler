package handlers

import "go/ast"

func handleValueSpec(spec ast.Spec) string {
	s := spec.(*ast.ValueSpec)
	code := ""
	code += handleValueSpecType(s.Type)
	code += " "
	code += handleValueSpecNames(s.Names)
	if s.Values != nil {
		code += " = "
		code += handleValueSpecValues(s.Values)
	}
	return code
}
