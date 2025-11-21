package transpile

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"

	"github.com/andygeiss/esp32-transpiler/transpile/handlers"
)

// Service specifies the api logic of transforming a source code format into another target format.
type Service interface {
	Start() error
}

const (
	// ErrorWorkerReaderIsNil ...
	ErrorWorkerReaderIsNil = "Reader should not be nil"
	// ErrorWorkerWriterIsNil ...
	ErrorWorkerWriterIsNil = "Writer should not be nil"
)

// defaultService specifies the api logic of transforming a source code format into another target format.
type defaultService struct {
	in  io.Reader
	out io.Writer
}

// NewService creates a a new transpile and returns its address.
func NewService(in io.Reader, out io.Writer) Service {
	return &defaultService{
		in:  in,
		out: out,
	}
}

// Start ...
func (s *defaultService) Start() error {
	if s.in == nil {
		return fmt.Errorf("Error: %s", ErrorWorkerReaderIsNil)
	}
	if s.out == nil {
		return fmt.Errorf("Error: %s", ErrorWorkerWriterIsNil)
	}

	// Read tokens from file by using Go's parser.
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "source.go", s.in, 0)
	if err != nil {
		return fmt.Errorf("ParseFile failed! %v", err)
	}

	// If source has no declarations then main it to an empty for loop.
	if file.Decls == nil {
		_, _ = fmt.Fprint(s.out, "void loop() {} void setup() {}")
		return nil
	}

	// Use Goroutines to work concurrently.
	count := len(file.Decls)
	done := make(chan bool, count)
	dst := make([]chan string, count)
	for i := 0; i < count; i++ {
		dst[i] = make(chan string, 1)
	}

	// Start a transpile with an individual channel for each declaration in the source file.
	for i, decl := range file.Decls {
		go handlers.HandleDecl(decl, dst[i], done)
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
			_, _ = s.out.Write([]byte(content))
		}
	}
	return nil
}
