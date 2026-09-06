package handlers

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

// handleBasicLit writes a literal the C++ compiler reads the same way Go does.
// Most come across as they stand; the ones that do not are Go's raw strings,
// its 0o octal and the underscores it allows inside a number.
func handleBasicLit(bl *ast.BasicLit) string {
	switch bl.Kind {
	case token.STRING:
		s, err := strconv.Unquote(bl.Value)
		if err != nil {
			unsupported(bl, "the string %s", bl.Value)
		}
		return strconv.Quote(s)
	case token.INT, token.FLOAT:
		return cNumber(bl)
	case token.CHAR:
		// A C++ char holds one byte; a rune above ASCII would not fit.
		r, err := strconv.Unquote(bl.Value)
		if err != nil || len(r) != 1 {
			unsupported(bl, "the character %s", bl.Value)
		}
		return bl.Value
	case token.IMAG:
		unsupported(bl, "a complex number")
	}
	return bl.Value
}

func cNumber(bl *ast.BasicLit) string {
	value := strings.ReplaceAll(bl.Value, "_", "")
	// Go writes octal 0o755 where C++ writes 0755.
	if len(value) > 2 && value[0] == '0' && (value[1] == 'o' || value[1] == 'O') {
		return "0" + value[2:]
	}
	return value
}
