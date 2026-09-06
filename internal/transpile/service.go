// Package transpile turns Go source code into an Arduino sketch.
package transpile

import (
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/andygeiss/esp32-transpiler/internal/transpile/handlers"
)

// Errors a caller branches on.
var (
	ErrNilReader = errors.New("reader is nil")
	ErrNilWriter = errors.New("writer is nil")
)

// emptySketch is what source that declares nothing the sketch needs becomes:
// the two functions the Arduino runtime always calls.
const emptySketch = "void setup() {}\nvoid loop() {}\n"

// defaultName is what diagnostics call the source when the caller names none.
const defaultName = "source.go"

// Service reads Go source and writes the matching Arduino sketch.
type Service struct {
	name string
	in   io.Reader
	out  io.Writer
}

// NewService creates a service that reads Go source from in and writes the
// sketch to out. name is what diagnostics call the source.
func NewService(name string, in io.Reader, out io.Writer) *Service {
	if name == "" {
		name = defaultName
	}
	return &Service{name: name, in: in, out: out}
}

// Start transpiles the whole source and writes the sketch in one go. A
// cancelled ctx stops it before it has written anything, and so does a
// construct outside the subset: half a sketch is worse than none.
func (s *Service) Start(ctx context.Context) error {
	if s.in == nil {
		return ErrNilReader
	}
	if s.out == nil {
		return ErrNilWriter
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, s.name, s.in, 0)
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}

	// Declaration by declaration, in source order, so the sketch reads the way
	// the Go file does.
	var sketch strings.Builder
	for _, decl := range file.Decls {
		if err := ctx.Err(); err != nil {
			return err
		}
		code, err := handlers.HandleDecl(fset, decl)
		if err != nil {
			return err
		}
		if code == "" {
			continue // an import that needs no header, say
		}
		sketch.WriteString(code)
		// One declaration a line: the sketch is opened in the Arduino IDE by
		// someone who has to read it.
		if !strings.HasSuffix(code, "\n") {
			sketch.WriteString("\n")
		}
	}

	// Source with nothing in it, or nothing that reaches the sketch, still has
	// to give the runtime the two functions it calls.
	if sketch.Len() == 0 {
		_, err = io.WriteString(s.out, emptySketch)
		return err
	}

	_, err = io.WriteString(s.out, sketch.String())
	return err
}
