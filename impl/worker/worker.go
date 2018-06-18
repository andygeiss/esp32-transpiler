package worker

import (
	"fmt"
	"github.com/andygeiss/esp32-transpiler/api/worker"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"
)

const (
	// ErrorWorkerReaderIsNil ...
	ErrorWorkerReaderIsNil = "Reader should not be nil"
	// ErrorWorkerWriterIsNil ...
	ErrorWorkerWriterIsNil = "Writer should not be nil"
)

var mapping worker.Mapping

// Worker specifies the api logic of transforming a source code format into another target format.
type Worker struct {
	in  io.Reader
	out io.Writer
}

// NewWorker creates a a new worker and returns its address.
func NewWorker(in io.Reader, out io.Writer, m worker.Mapping) worker.Worker {
	mapping = m
	return &Worker{in, out}
}

// Start ...
func (w *Worker) Start() error {
	if w.in == nil {
		return fmt.Errorf("Error: %s", ErrorWorkerReaderIsNil)
	}
	if w.out == nil {
		return fmt.Errorf("Error: %s", ErrorWorkerWriterIsNil)
	}
	// Read tokens from file by using Go's parser.
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "source.go", w.in, 0)
	if err != nil {
		return fmt.Errorf("ParseFile failed! %v", err)
	}
	// If source has no declarations then main it to an empty for loop.
	if file.Decls == nil {
		fmt.Fprint(w.out, "void loop() {} void setup() {}")
		return nil
	}
	// Use Goroutines to work concurrently.
	count := len(file.Decls)
	done := make(chan bool, count)
	dst := make([]chan string, count)
	for i := 0; i < count; i++ {
		dst[i] = make(chan string, 1)
	}
	// Start a worker with an individual channel for each declaration in the source file.
	for i, decl := range file.Decls {
		go handleDecl(i, decl, dst[i], done)
	}
	// Wait for all workers are done.
	for i := 0; i < count; i++ {
		select {
		case <-done:
		}
	}
	// Print the ordered result.
	for i := 0; i < count; i++ {
		for content := range dst[i] {
			w.out.Write([]byte(content))
		}
	}
	// Print the AST.
	// ast.Fprint(os.Stderr, fset, file, nil)
	return nil
}

func handleAssignStmt(as *ast.AssignStmt) string {
	code := handleAssignStmtExpr(as.Lhs)
	code += as.Tok.String()
	code += handleAssignStmtExpr(as.Rhs)
	return code
}

func handleAssignStmtExpr(e []ast.Expr) string {
	ops := make([]string, 0)
	code := ""
	for _, op := range e {
		ops = append(ops, handleExpr(op))
	}
	code += strings.Join(ops, ",")
	return code
}

func handleBasicLit(bl *ast.BasicLit) string {
	return bl.Value
}

func handleBinaryExpr(expr ast.Expr) string {
	be := expr.(*ast.BinaryExpr)
	code := handleExpr(be.X)
	code += be.Op.String()
	code += handleExpr(be.Y)
	return code
}

func handleCallExpr(expr *ast.CallExpr) string {
	code := handleExpr(expr.Fun)
	code += "("
	args := make([]string, 0)
	for _, arg := range expr.Args {
		args = append(args, handleExpr(arg))
	}
	code += strings.Join(args, ",")
	code += ")"
	return code
}

func handleDecl(id int, decl ast.Decl, dst chan<- string, done chan<- bool) {
	code := ""
	switch d := decl.(type) {
	case *ast.FuncDecl:
		code += handleFuncDecl(d)
	case *ast.GenDecl:
		code += handleGenDecl(d)
	}
	dst <- code
	close(dst)
	done <- true
}

func handleDeclStmt(stmt *ast.DeclStmt) string {
	code := ""
	switch decl := stmt.Decl.(type) {
	case *ast.GenDecl:
		code += handleGenDecl(decl)
	}
	return code
}

func handleExpr(expr ast.Expr) string {
	code := ""
	switch e := expr.(type) {
	case *ast.BasicLit:
		code += handleBasicLit(e)
	case *ast.BinaryExpr:
		code += handleBinaryExpr(e)
	case *ast.CallExpr:
		code += handleCallExpr(e)
	case *ast.Ident:
		code += handleIdent(e)
	case *ast.SelectorExpr:
		code += handleSelectorExpr(e)
	}
	return code
}

func handleExprStmt(stmt *ast.ExprStmt) string {
	code := ""
	switch x := stmt.X.(type) {
	case *ast.CallExpr:
		code += handleCallExpr(x)
	}
	return code
}

func handleFuncDecl(decl ast.Decl) string {
	fd := decl.(*ast.FuncDecl)
	code := ""
	name := ""
	code += handleFuncDeclType(fd.Type)
	code += " "
	name = handleFuncDeclName(fd.Name)
	if name == "NewController" {
		return ""
	}
	code += name
	code += "("
	code += handleFuncDeclParams(fd.Type)
	code += ") {"
	code += handleBlockStmt(fd.Body)
	code += "}"
	return code
}

func handleFuncDeclParams(t *ast.FuncType) string {
	code := ""
	if t.Params == nil || t.Params.List == nil {
		return code
	}
	values := make([]string, 0)
	for _, field := range t.Params.List {
		ftype := ""
		switch ft := field.Type.(type) {
		case *ast.Ident:
			ftype = handleIdent(ft)
		}
		for _, names := range field.Names {
			values = append(values, ftype+" "+names.Name)
		}
	}
	code += strings.Join(values, ",")
	return code
}

func handleBlockStmt(body *ast.BlockStmt) string {
	code := ""
	if body == nil {
		return code
	}
	for _, stmt := range body.List {
		code += handleStmt(stmt)
	}
	return code
}

func handleBranchStmt(stmt *ast.BranchStmt) string {
	return "break;"
}

func handleCaseClause(cc *ast.CaseClause) string {
	code := "case "
	clauses := make([]string, 0)
	for _, clause := range cc.List {
		clauses = append(clauses, handleExpr(clause))
	}
	code += strings.Join(clauses, ",")
	code += ":"
	for _, body := range cc.Body {
		code += handleStmt(body)
	}
	return code
}

func handleFuncDeclName(ident *ast.Ident) string {
	code := ""
	if ident == nil {
		return code
	}
	code += ident.Name
	code = mapping.Apply(code)
	return code
}

func handleFuncDeclType(t *ast.FuncType) string {
	code := ""
	if t.Results == nil {
		code = "void"
	}
	return code
}

func handleGenDecl(decl ast.Decl) string {
	gd := decl.(*ast.GenDecl)
	code := ""
	switch gd.Tok {
	case token.CONST:
		code += "const "
	case token.VAR:
		code += ""
	}
	code += handleSpecs(gd.Specs)
	return code
}

func handleIdent(expr ast.Expr) string {
	ident := expr.(*ast.Ident)
	code := ""
	switch ident.Name {
	case "string":
		code += "char*"
	default:
		code += ident.Name
	}
	return code
}

func handleIfStmt(stmt *ast.IfStmt) string {
	cond := handleExpr(stmt.Cond)
	body := handleBlockStmt(stmt.Body)
	code := fmt.Sprintf(`if (%s) { %s }`, cond, body)
	if stmt.Else != nil {
		tail := handleBlockStmt(stmt.Else.(*ast.BlockStmt))
		code += fmt.Sprintf(" else { %s }", tail)
	}
	return code
}

func handleImportSpec(spec ast.Spec) string {
	s := spec.(*ast.ImportSpec)
	code := ""
	if s.Name != nil {
		name := handleIdent(s.Name)
		name = mapping.Apply(name)
		if name != "" {
			if name != "controller" {
				code = "#include <" + name + ".h>\n"
			}
		}
	}
	return code
}

func handleSelectorExpr(expr ast.Expr) string {
	s := expr.(*ast.SelectorExpr)
	code := ""
	switch x := s.X.(type) {
	case *ast.Ident:
		code += handleIdent(x)
	}
	code += "."
	code += handleIdent(s.Sel)
	code = mapping.Apply(code)
	return code
}

func handleSpecs(specs []ast.Spec) string {
	code := ""
	for _, spec := range specs {
		switch spec.(type) {
		case *ast.ImportSpec:
			code += handleImportSpec(spec)
		case *ast.ValueSpec:
			code += handleValueSpec(spec) + ";"
		}
	}
	return code
}

func handleStmt(stmt ast.Stmt) string {
	code := ""
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		code += handleAssignStmt(s)
		code += ";"
	case *ast.BranchStmt:
		code += handleBranchStmt(s)
	case *ast.CaseClause:
		code += handleCaseClause(s)
	case *ast.DeclStmt:
		code += handleDeclStmt(s)
	case *ast.ExprStmt:
		code += handleExprStmt(s)
		code += ";"
	case *ast.ForStmt:
		code += handleForStmt(s)
	case *ast.IfStmt:
		code += handleIfStmt(s)
	case *ast.SwitchStmt:
		code += handleSwitchStmt(s)
	}
	return code
}

func handleForStmt(stmt *ast.ForStmt) string {
	code := ""
	if stmt.Init == nil && stmt.Post == nil {
		code += "while"
	} else {
		code += "for"
	}
	code += "(" // stmt.Init
	code += handleBinaryExpr(stmt.Cond) // stmt.Cond
	code += "" // stmt.Post
	code += ") {"
	code += handleBlockStmt(stmt.Body) // stmt.Body
	code += "}"
	return code
}

func handleSwitchStmt(stmt *ast.SwitchStmt) string {
	code := "switch ("
	code += handleExpr(stmt.Tag)
	code += "){"
	code += handleBlockStmt(stmt.Body)
	code += "}"
	return code
}

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

func handleValueSpecNames(names []*ast.Ident) string {
	code := ""
	for _, name := range names {
		code += handleIdent(name)
	}
	return code
}

func handleValueSpecType(expr ast.Expr) string {
	code := ""
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		code += handleSelectorExpr(t)
	case *ast.Ident:
		code += handleIdent(t)
	}
	return code
}

func handleValueSpecValues(values []ast.Expr) string {
	code := ""
	for _, value := range values {
		switch v := value.(type) {
		case *ast.BasicLit:
			code += handleBasicLit(v)
		}
	}
	return code
}
