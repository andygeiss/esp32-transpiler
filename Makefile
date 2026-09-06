# Copied from the baseline (stack/makefile.md). Adjust per its rule 5; record
# any other deviation in the README.

# The main package: this module ships one binary, so it sits at the root.
MAIN = .

# Targets are alphabetical, so the default is named rather than first.
.DEFAULT_GOAL = check
.PHONY: build check ci clean fmt run test

# Release-shaped local binary in bin/ (go build creates the directory).
build:
	CGO_ENABLED=0 go build -trimpath -o bin/ $(MAIN)

# Default. Every gate, in this order (operations/ci.md), against the working
# tree. Run before every commit.
check:
	test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	go vet ./...
	go fix -diff ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go mod tidy -diff
	go test -race -shuffle=on ./...
	CGO_ENABLED=0 go build -trimpath ./...

# The same gates against the commit: a file never added, or a .env, cannot
# make it green. Run before every push. go version runs first, inside the
# copy, so the run records which toolchain ran. The archive goes through a
# file so git's exit status stops the run; one shell line so the trap cleans
# up however check ends.
ci:
	t=$$(mktemp); d=$$(mktemp -d); trap 'rm -rf "$$t" "$$d"' EXIT; git archive -o "$$t" HEAD && tar -xf "$$t" -C "$$d" && go -C "$$d" version && $(MAKE) -C "$$d" check

clean:
	rm -rf bin/ esp32-transpiler

# goimports first: go fix type-checks, so a missing import would stop the
# recipe before goimports could add it. go fix manages the imports its own
# rewrites need.
fmt:
	go run golang.org/x/tools/cmd/goimports@latest -w .
	go fix ./...

# Loads .env when it is there, so a local start is one command. Only run:
# check and test MUST NOT depend on a developer's machine (rule 6). One shell
# line, because each recipe line gets its own shell.
run:
	set -a; if [ -f .env ]; then . ./.env; fi; set +a; go run $(MAIN)

# The inner loop.
test:
	go test -race -shuffle=on ./...
