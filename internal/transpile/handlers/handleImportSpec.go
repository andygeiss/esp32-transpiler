package handlers

import "go/ast"

func handleImportSpec(spec ast.Spec) string {
	s := spec.(*ast.ImportSpec)
	code := ""
	if s.Name != nil {
		name := handleIdent(s.Name)
		if val, ok := mapping[name]; ok {
			name = val
		}
		if name != "" {
			if name != "controller" {
				code = "#include <" + name + ".h>\n"
			}
		}
	}
	return code
}
