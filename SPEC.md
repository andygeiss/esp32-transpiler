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
proves it, so changing either changes what the tool promises. Everything outside
that part is refused by name and line, and no sketch is written: a construct the
tool passes over silently would reach the board as a sketch that compiles and
behaves differently from the Go, and finding that out on the board is what this
tool exists to avoid. No baseline rule is waived; the README says so.

A release number answers one question: which
[esp32-controller](https://github.com/andygeiss/esp32-controller) release this
transpiler is built to fit.

- **Major and minor mirror the controller.** `v0.3.0` fits esp32-controller
  `v0.3.0`, and fitting means every Arduino call that version exports has a row
  in `mapping.go`.
- **The patch is the transpiler's own.** A fix that does not change which
  controller release the tool fits takes the next patch on the current line, the
  way `v0.2.1` did between the controller's `v0.2.0` and its `v0.3.0`.
- **The number carries no semver, so the notes must.** A change to a flag, an
  exit code, or the meaning of a generated sketch is a break whatever the number
  does that release. Say so in the tag message and in the release notes, in those
  words. This is the trade `baseline-reference` makes, and the cost is the same
  one: the version says nothing about the tool's own contract.

**Done means:** The CLI tool checklist is walked, and `make ci` is green on the
commit.
