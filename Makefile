.DEFAULT_GOAL := all

.PHONY: all test vet cover

all: test vet

test:
	go test ./...

vet:
	go vet ./...

cover:
	go test -cover ./...
