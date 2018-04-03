TS=$(shell date -u '+%Y/%m/%d %H:%M:%S')

all: test

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
