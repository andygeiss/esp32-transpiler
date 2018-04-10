TS=$(shell date -u '+%Y/%m/%d %H:%M:%S')

all: install test

install:
	@echo $(TS) Installing...
	@go get -u github.com/andygeiss/assert
	@go get -u github.com/andygeiss/esp32-controller
	@go get -u github.com/andygeiss/esp32-transpiler
	@echo $(TS) Done.

test:
	@echo $(TS) Testing ...
	@go test -v ./...
	@echo $(TS) Done.
