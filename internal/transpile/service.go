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

// emptySketch is what source without declarations becomes: the two functions
// the Arduino runtime always calls.
const emptySketch = "void loop() {} void setup() {}"

// Service reads Go source and writes the matching Arduino sketch.
type Service struct {
	in  io.Reader
	out io.Writer
}

// NewService creates a service that reads Go source from in and writes the
// sketch to out.
func NewService(in io.Reader, out io.Writer) *Service {
	return &Service{in: in, out: out}
}

// Start transpiles the whole source and writes the sketch in one go. A
// cancelled ctx stops it before it has written anything.
func (s *Service) Start(ctx context.Context) error {
	if s.in == nil {
		return ErrNilReader
	}
	if s.out == nil {
		return ErrNilWriter
	}

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "source.go", s.in, 0)
	if err != nil {
		return fmt.Errorf("parsing source: %w", err)
	}

	if len(file.Decls) == 0 {
		_, err := io.WriteString(s.out, emptySketch)
		return err
	}

	// Declaration by declaration, in source order, so the sketch reads the way
	// the Go file does.
	var sketch strings.Builder
	for _, decl := range file.Decls {
		if err := ctx.Err(); err != nil {
			return err
		}
		sketch.WriteString(handlers.HandleDecl(decl))
	}

	_, err = io.WriteString(s.out, sketch.String())
	return err
}
