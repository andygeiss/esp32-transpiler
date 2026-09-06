package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// One Go source and the whole sketch it transpiles to.
const (
	goSource = "package test\nfunc foo() {}\n"
	sketch   = "void foo() {}"
)

// writeSource puts a Go source file in a fresh directory and returns its path.
func writeSource(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "source.go")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing the source: %v", err)
	}
	return path
}

func TestRun_WritesTheSketchToAFile(t *testing.T) {
	t.Parallel()
	source := writeSource(t, goSource)
	dir := filepath.Dir(source)
	target := filepath.Join(dir, "sketch.ino")
	var stdout, stderr bytes.Buffer

	if err := run(t.Context(), []string{"-source", source, "-target", target}, &stdout, &stderr); err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading the sketch: %v", err)
	}
	if string(got) != sketch {
		t.Errorf("got %q, want %q", got, sketch)
	}
	if stdout.Len() != 0 {
		t.Errorf("got %q on stdout, want nothing", stdout.String())
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("reading the sketch: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o644 {
		t.Errorf("got mode %v, want -rw-r--r--", perm)
	}
	// The file the atomic write renamed onto the target must be gone.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading the directory: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("got %d files, want 2 (the source and the sketch)", len(entries))
	}
}

func TestRun_TargetDashWritesTheSketchToStdout(t *testing.T) {
	t.Parallel()
	source := writeSource(t, goSource)
	var stdout, stderr bytes.Buffer

	if err := run(t.Context(), []string{"-source", source, "-target", "-"}, &stdout, &stderr); err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}

	if stdout.String() != sketch {
		t.Errorf("got %q, want %q", stdout.String(), sketch)
	}
}

func TestRun_UsageErrors(t *testing.T) {
	t.Parallel()
	source := writeSource(t, goSource)
	tests := []struct {
		name string
		args []string
	}{
		{"no flags at all", nil},
		{"a source without a target", []string{"-source", source}},
		{"a target without a source", []string{"-target", "-"}},
		{"a flag the tool does not have", []string{"-nope"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := run(t.Context(), tt.args, &stdout, &stderr)
			if !errors.Is(err, errUsage) {
				t.Errorf("got %v, want %v", err, errUsage)
			}
			if stdout.Len() != 0 {
				t.Errorf("got %q on stdout, want nothing", stdout.String())
			}
			if stderr.Len() == 0 {
				t.Error("got nothing on stderr, want the usage text")
			}
		})
	}
}

func TestRun_HelpExitsWithoutAnError(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer

	if err := run(t.Context(), []string{"-h"}, &stdout, &stderr); err != nil {
		t.Errorf("got %v, want no error", err)
	}

	if !strings.Contains(stderr.String(), "-source") {
		t.Errorf("got %q on stderr, want the usage text", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("got %q on stdout, want nothing", stdout.String())
	}
}

func TestRun_VersionGoesToStdout(t *testing.T) {
	t.Parallel()
	var stdout, stderr bytes.Buffer

	if err := run(t.Context(), []string{"-version"}, &stdout, &stderr); err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}

	if strings.TrimSpace(stdout.String()) == "" {
		t.Error("got nothing on stdout, want a version")
	}
}

func TestRun_AMissingSourceIsNotAUsageError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	var stdout, stderr bytes.Buffer

	err := run(t.Context(),
		[]string{"-source", filepath.Join(dir, "gone.go"), "-target", filepath.Join(dir, "sketch.ino")},
		&stdout, &stderr)

	if err == nil {
		t.Fatal("got no error, want one")
	}
	if errors.Is(err, errUsage) {
		t.Errorf("got %v, want a plain error so the tool exits 1", err)
	}
}

func TestRun_BrokenSourceLeavesTheOldSketchIntact(t *testing.T) {
	t.Parallel()
	source := writeSource(t, "package test\nfunc foo( {}")
	target := filepath.Join(filepath.Dir(source), "sketch.ino")
	const old = "void old() {}"
	if err := os.WriteFile(target, []byte(old), 0o644); err != nil {
		t.Fatalf("writing the old sketch: %v", err)
	}

	var stdout, stderr bytes.Buffer
	if err := run(t.Context(), []string{"-source", source, "-target", target}, &stdout, &stderr); err == nil {
		t.Fatal("got no error, want a parse error")
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("reading the sketch: %v", err)
	}
	if string(got) != old {
		t.Errorf("got %q, want the old sketch %q", got, old)
	}
}
