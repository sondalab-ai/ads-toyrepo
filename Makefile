.DEFAULT_GOAL := all

.PHONY: all test vet

all: test vet

test:
	go test ./...

vet:
	go vet ./...
