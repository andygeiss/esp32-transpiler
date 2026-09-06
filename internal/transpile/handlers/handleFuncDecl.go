package handlers

import "go/ast"

func handleFuncDecl(fd *ast.FuncDecl) string {
	if fd.Recv != nil {
		unsupported(fd, "a method")
	}
	// A Go function without a body is written somewhere else, in assembly.
	// Giving the sketch an empty one would quietly answer every call with
	// nothing.
	if fd.Body == nil {
		unsupported(fd, "a function without a body")
	}

	code := ""
	body := fd.Body
	if _, mapped := mapping[fd.Name.Name]; mapped {
		// The Arduino runtime calls setup and loop itself: it passes them
		// nothing and has nowhere to put a result. Their mapped names carry
		// the void already.
		rejectRuntimeSignature(fd)
		rejectValueReturns(fd)
		code = handleFuncDeclName(fd.Name)
		body = trimTrailingReturn(body)
	} else {
		result := handleFuncDeclType(fd.Type)
		code = result + " " + handleFuncDeclName(fd.Name)
		if result == "void" {
			body = trimTrailingReturn(body)
		}
	}
	code += "("
	code += handleFuncDeclParams(fd.Type)
	code += ") {"
	code += handleBlockStmt(body)
	code += "}"
	return code
}

// trimTrailingReturn drops a return that ends a void function, where the
// closing brace says the same thing. One anywhere else means stop now, and a
// sketch that carried on instead would do what the Go never did.
func trimTrailingReturn(body *ast.BlockStmt) *ast.BlockStmt {
	if len(body.List) == 0 {
		return body
	}
	last, ok := body.List[len(body.List)-1].(*ast.ReturnStmt)
	if !ok || !endsFunction(last) {
		return body
	}
	trimmed := *body
	trimmed.List = body.List[:len(body.List)-1]
	return &trimmed
}

// rejectRuntimeSignature keeps Setup and Loop to the shape the board calls. A
// parameter or a result other than error would leave a sketch that compiles
// and then fails to link.
func rejectRuntimeSignature(fd *ast.FuncDecl) {
	if p := fd.Type.Params; p != nil && len(p.List) > 0 {
		unsupported(fd, "%s with a parameter", fd.Name.Name)
	}
	r := fd.Type.Results
	if r == nil || len(r.List) == 0 {
		return
	}
	if len(r.List) == 1 && len(r.List[0].Names) == 0 {
		if ident, ok := r.List[0].Type.(*ast.Ident); ok && ident.Name == "error" {
			return
		}
	}
	unsupported(fd, "%s returning anything but error", fd.Name.Name)
}

// rejectValueReturns stops a setup or loop from handing something back. The
// board has nowhere to put it, and C++ will not return a value out of a void
// function, so the sketch would not compile.
func rejectValueReturns(fd *ast.FuncDecl) {
	if fd.Body == nil {
		return
	}
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false // its returns are its own, and it is turned down later
		}
		ret, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		if !endsFunction(ret) {
			unsupported(ret, "%s returning a value", fd.Name.Name)
		}
		return true
	})
}
