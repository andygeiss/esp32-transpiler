package ino

import (
	"fmt"
	"github.com/andygeiss/esp32-transpiler/business/worker"
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

// Worker specifies the business logic of transforming a source code format into another target format.
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
	//ast.Fprint(os.Stderr, fset, file, nil)
	return nil
}

func handleAssignStmt(stmt ast.Stmt) string {
	a := stmt.(*ast.AssignStmt)
	code := ""
	code += handleAssignStmtExpr(a.Lhs)
	code += a.Tok.String()
	code += handleAssignStmtExpr(a.Rhs)
	return code
}

func handleAssignStmtExpr(expr []ast.Expr) string {
	code := ""
	ops := make([]string, 0)
	for _, op := range expr {
		switch o := op.(type) {
		case *ast.BasicLit:
			ops = append(ops, handleBasicLit(o))
		case *ast.BinaryExpr:
			ops = append(ops, handleBinaryExpr(o))
		case *ast.CallExpr:
			ops = append(ops, handleCallExpr(o))
		case *ast.Ident:
			ops = append(ops, handleIdent(o))
		}
	}
	code += strings.Join(ops, ",")
	return code
}

func handleBasicLit(expr ast.Expr) string {
	bl := expr.(*ast.BasicLit)
	code := ""
	code += bl.Value
	return code
}

func handleBinaryExpr(expr ast.Expr) string {
	e := expr.(*ast.BinaryExpr)
	code := ""
	switch x := e.X.(type) {
	case *ast.Ident:
		code += handleIdent(x)
	}
	code += e.Op.String()
	switch y := e.Y.(type) {
	case *ast.Ident:
		code += handleIdent(y)
	}
	return code
}

func handleCallExpr(expr *ast.CallExpr) string {
	code := ""
	switch e := expr.Fun.(type) {
	case *ast.Ident:
		code += handleIdent(e)
	case *ast.SelectorExpr:
		code += handleSelectorExpr(e)
	}
	code += "("
	args := make([]string, 0)
	for _, arg := range expr.Args {
		switch a := arg.(type) {
		case *ast.BasicLit:
			args = append(args, handleBasicLit(a))
		case *ast.CallExpr:
			args = append(args, handleCallExpr(a))
		case *ast.SelectorExpr:
			args = append(args, handleSelectorExpr(a))
		}
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
	code += handleFuncDeclBody(fd.Body)
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

func handleFuncDeclBody(body *ast.BlockStmt) string {
	code := ""
	if body == nil {
		return code
	}
	for _, stmt := range body.List {
		switch s := stmt.(type) {
		case *ast.AssignStmt:
			code += handleAssignStmt(s)
			code += ";"
		case *ast.DeclStmt:
			code += handleDeclStmt(s)
		case *ast.ExprStmt:
			code += handleExprStmt(s)
			code += ";"
		}
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

func handleValueSpec(spec ast.Spec) string {
	s := spec.(*ast.ValueSpec)
	code := ""
	code += handleValueSpecType(s.Type)
	code += " "
	code += handleValueSpecNames(s.Names)
	code += " = "
	code += handleValueSpecValues(s.Values)
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
