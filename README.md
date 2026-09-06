# ESP32 Transpiler

[![License](https://img.shields.io/github/license/andygeiss/esp32-transpiler)](https://github.com/andygeiss/esp32-transpiler/blob/master/LICENSE)
[![Releases](https://img.shields.io/github/v/release/andygeiss/esp32-transpiler)](https://github.com/andygeiss/esp32-transpiler/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/esp32-transpiler)](https://goreportcard.com/report/github.com/andygeiss/esp32-transpiler)

`esp32-transpiler` turns a Go source file into an Arduino sketch for the ESP32.
It is for people who would rather write their controller logic in Go and check
it with `go test` than flash a board to find out whether it works.

## Install

```sh
go install github.com/andygeiss/esp32-transpiler@latest
```

Version numbers track [esp32-controller](https://github.com/andygeiss/esp32-controller):
the same number on both means this transpiler covers every Arduino call that
version of the controller exports.

## Transpile a controller

Write the controller in Go. `Setup` and `Loop` become the two functions the
Arduino runtime calls, and the packages you import come from
[esp32-controller](https://github.com/andygeiss/esp32-controller), which mirrors
the Arduino calls in Go so the same code runs under `go test`:

```go
package controller

import "github.com/andygeiss/esp32-controller/serial"

// Setup runs once when the board powers up.
func Setup() error {
	serial.Begin(serial.BaudRate115200)
	return nil
}

// Loop runs over and over, as fast as it returns.
func Loop() error {
	serial.Println("hello")
	return nil
}
```

Turn it into a sketch:

```sh
esp32-transpiler -source controller.go -target controller.ino
```

`controller.ino` now holds this, ready for the ESP32 toolchain:

```c
void setup() {Serial.begin(115200);}
void loop() {Serial.println("hello");}
```

## Options

```
-source <file>  Go source to read; "-" reads standard input
-target <file>  Arduino sketch to write; "-" writes standard output
-version        print the version and exit
-h              print this help and exit
```

The sketch is the only thing that reaches standard output, so `-target -` pipes.
Everything else — usage, errors — goes to standard error. The tool exits 0 when
it wrote the sketch, 1 when it could not, and 2 when the command line was wrong.

A run that fails leaves the sketch that was already there untouched: the new one
is written beside it and renamed onto it, and a rename either happens or does
not.

## What it transpiles

Only a small part of the [Go language specification](https://go.dev/ref/spec).
[`mapping.go`](internal/transpile/handlers/mapping.go) is the full list of names
it rewrites, and [`service_test.go`](internal/transpile/service_test.go) shows
every construct it handles.

An import only becomes an `#include` when you give it a name:
`import wifi "…/wifi"` writes `#include <WiFi.h>`, while a plain import writes
nothing. Serial, pins and timers need no header on the ESP32, so import those
without a name.

Anything outside that part stops the run, with the construct and the line it
sits on:

```sh
$ esp32-transpiler -source controller.go -target controller.ino
esp32-transpiler: controller.go:12:3: a range loop is not supported
```

A sketch is not written, so the one already there is left alone. The tool
refuses rather than skips because the alternative is a sketch that compiles and
does something other than the Go said — which is the thing you ran `go test` to
rule out.

Three things it will not do:

- **It does not manage memory.** Go collects garbage on its own; C++ on a
  microcontroller does not, and the transpiler adds nothing to make up for it.
- **Go strings become `const char*`,** which a sketch can keep on the stack.
  Joining two of them with `+` is Go, not C++, so the tool turns it down when it
  can see a literal on either side. It cannot see one between two variables:
  give the sketch `strcat` or an Arduino `String` instead.
- **It does not type-check.** `go build` and `go test` do that, and running them
  on the controller first is the point of writing it in Go. One difference it
  cannot see for you: a Go `int` is 64 bits on the machine the tests run on and
  32 on the ESP32, so arithmetic that only just fits in Go can overflow on the
  board. Size the type yourself where it is close.

## Build and test it

```sh
make        # every gate: format, vet, fix, staticcheck, govulncheck, tidy, test, build
make test   # the inner loop
make build  # a release-shaped binary in bin/
```

`make ci` runs those same gates against the commit rather than the working tree.
It is what has to be green before a push — there is no CI server.

## Baseline

This repository follows the [engineering baseline](https://github.com/andygeiss/baseline),
CLI tool track. No baseline rule is waived.

## License

MIT — see [LICENSE](LICENSE).
