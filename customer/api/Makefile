.PHONY: fmt install test

REVISION := $(shell git rev-parse HEAD)

default: install

fmt:
	go fmt ./...

install:
	./scripts/install.sh

test:
	go test ./...
