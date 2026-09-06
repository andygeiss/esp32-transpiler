package handlers

import "go/ast"

// handleImportSpec turns a named import into the header the Arduino needs.
// A plain import names nothing to include: serial, pins and timers are in the
// ESP32 core already. The controller package itself is Go-only scaffolding.
func handleImportSpec(s *ast.ImportSpec) string {
	if s.Name == nil {
		return ""
	}
	name := s.Name.Name
	switch name {
	case "_", ".", "controller":
		return ""
	}
	if val, ok := mapping[name]; ok {
		name = val
	}
	return "#include <" + name + ".h>\n"
}
