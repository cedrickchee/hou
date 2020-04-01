.PHONY: build test deps clean

all: build
	@./hou

deps:
	@go get ./...

build:
	@go build .

test:
	@go test -v -cover -coverprofile=coverage.out -covermode=atomic ./...

clean:
	@rm -rf hou