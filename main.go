// Command esp32-transpiler turns a Go source file into an Arduino sketch.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/andygeiss/esp32-transpiler/internal/transpile"
)

const usage = `esp32-transpiler turns a Go source file into an Arduino sketch.

Usage:
    esp32-transpiler -source <file> -target <file>

Options:
    -source <file>  Go source to read; "-" reads standard input
    -target <file>  Arduino sketch to write; "-" writes standard output
    -version        print the version and exit
    -h              print this help and exit

Example:
    esp32-transpiler -source blink.go -target blink.ino
`

// errUsage means the command line was wrong, and its message is already on
// stderr. main turns it into exit code 2.
var errUsage = errors.New("usage error")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	switch {
	case err == nil:
	case errors.Is(err, errUsage):
		os.Exit(2)
	default:
		fmt.Fprintf(os.Stderr, "esp32-transpiler: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("esp32-transpiler", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }
	source := fs.String("source", "", `Go source to read ("-" reads standard input)`)
	target := fs.String("target", "", `Arduino sketch to write ("-" writes standard output)`)
	showVersion := fs.Bool("version", false, "print the version and exit")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil // -h: fs.Usage already printed it
		}
		return errUsage
	}

	if *showVersion {
		fmt.Fprintln(stdout, version())
		return nil
	}
	if *source == "" || *target == "" {
		fs.Usage()
		return errUsage
	}
	// Writing the sketch over the Go it came from would leave nothing to
	// transpile the next time.
	if *source != "-" && *target != "-" && sameFile(*source, *target) {
		fmt.Fprintln(stderr, "esp32-transpiler: -source and -target name the same file")
		return errUsage
	}

	in, err := open(*source)
	if err != nil {
		return err
	}
	defer in.Close()

	var sketch bytes.Buffer
	// The error already names the source and the line it stopped on, the way
	// a compiler does, so it goes out as it is.
	if err := transpile.NewService(sourceName(*source), in, &sketch).Start(ctx); err != nil {
		return err
	}

	// Nothing is written after a Ctrl-C, so an interrupted run leaves the old
	// sketch where it was.
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("interrupted before writing %s: %w", *target, err)
	}
	return write(*target, sketch.Bytes(), stdout)
}

// sourceName is what diagnostics call the source. Standard input has no name
// of its own.
func sourceName(source string) string {
	if source == "-" {
		return "<stdin>"
	}
	return source
}

// sameFile reports whether two names reach one file, whatever route they take
// to it. A target that is not there yet cannot be the source.
func sameFile(a, b string) bool {
	ai, err := os.Stat(a)
	if err != nil {
		return false
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false
	}
	return os.SameFile(ai, bi)
}

// open returns the source. Standard input has no file to close, so it gets a
// closer that does nothing.
func open(name string) (io.ReadCloser, error) {
	if name == "-" {
		return io.NopCloser(os.Stdin), nil
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("opening source: %w", err)
	}
	return f, nil
}

// write puts the sketch where -target asks. A file is written next to the
// target and renamed onto it, because a rename is atomic: a run that dies
// halfway leaves the previous sketch intact rather than a truncated one.
func write(name string, sketch []byte, stdout io.Writer) error {
	if name == "-" {
		_, err := stdout.Write(sketch)
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(name), filepath.Base(name)+".tmp")
	if err != nil {
		return fmt.Errorf("creating temporary sketch: %w", err)
	}
	defer os.Remove(tmp.Name()) // does nothing once the rename below succeeded
	if _, err := tmp.Write(sketch); err != nil {
		tmp.Close()
		return fmt.Errorf("writing %s: %w", tmp.Name(), err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", tmp.Name(), err)
	}
	// CreateTemp makes the file readable by its owner only; a sketch is not a
	// secret, and the Arduino IDE may run as someone else.
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return fmt.Errorf("setting permissions on %s: %w", tmp.Name(), err)
	}
	if err := os.Rename(tmp.Name(), name); err != nil {
		return fmt.Errorf("renaming %s to %s: %w", tmp.Name(), name, err)
	}
	return nil
}

// version reports what the toolchain stamped: the tag for a go install, a
// version derived from git for a build out of a checkout.
func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	// A build out of a checkout has no tag, so the commit it came from is the
	// next best answer.
	var revision, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value
		}
	}
	if revision == "" {
		return "unknown" // no VCS metadata to fall back on
	}
	if len(revision) > 12 {
		revision = revision[:12]
	}
	if modified == "true" {
		return "devel-" + revision + "-dirty"
	}
	return "devel-" + revision
}
