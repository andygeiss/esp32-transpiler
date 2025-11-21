set dotenv-load

# Compile the Go sources.
build: test
    @go build \
    -ldflags "-s -w" \
    -o ./bin/$(basename $PWD) ./main.go

# Install the binary to the $HOME/bin directory.
install: build
    @cp ./bin/$(basename $PWD) $HOME/bin/$(basename $PWD)

# Run the Go sources.
run: build
    @./bin/$(basename $PWD)

# Test the Go sources (Units).
test:
    @GOTOOLCHAIN=go1.25.4+auto go test -v -coverprofile=.coverprofile.out ./...
