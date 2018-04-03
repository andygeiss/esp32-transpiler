APPNAME=$(shell basename `pwd`)
LDFLAGS="-s"
TS=$(shell date -u '+%Y/%m/%d %H:%M:%S')

all: clean test build install

build/$(APPNAME):
	@echo $(TS) Building $(APPAME) ...
	@go build -ldflags $(LDFLAGS) -o build/$(APPNAME) main.go
	@echo $(TS) Done.

build: build/$(APPNAME)

clean:
	@echo $(TS) Cleaning up previous build ...
	@rm -f build/*
	@echo $(TS) Done.

install:
	@echo $(TS) Installing $(APPNAME) ...
	@cp build/$(APPNAME) $(GOPATH)/bin/
	@mkdir -p $(HOME)/esp32/
	@cp mapping.json $(HOME)/esp32/mapping.json
	@echo $(TS) Done.

packages:
	@echo $(TS) Installing Go packages ...
	@go get -u github.com/andygeiss/assert
	@go get -u github.com/andygeiss/esp32-controller
	@go get -u github.com/andygeiss/log
	@echo $(TS) Done.

test:
	@echo $(TS) Testing ...
	@go test -v ./...
	@echo $(TS) Done.
