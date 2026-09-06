# SPEC

**Job:** Turn a Go source file into an Arduino sketch the ESP32 toolchain can
compile, in one command.

**Why:** Controller logic written in the Arduino IDE has to be compiled and
flashed onto a board before anyone finds out whether it works. The same logic
written in Go is answered by `go test` in a second, and this tool turns the
tested Go into the sketch that ships.

**Guardrails:** The tool covers a small part of the Go language on purpose.
[`mapping.go`](internal/transpile/handlers/mapping.go) is the whole Go-to-Arduino
contract and [`service_test.go`](internal/transpile/service_test.go) is what
proves it, so changing either changes what the tool promises — that is a major
release. No baseline rule is waived; the README says so.

**Done means:** The CLI tool checklist is walked, and `make ci` is green on the
commit.
