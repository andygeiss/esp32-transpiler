TS=$(shell date -u '+%Y/%m/%d %H:%M:%S')

all: test

test:
	@echo $(TS) Testing ...
	@go test -v ./...
	@echo $(TS) Done.
