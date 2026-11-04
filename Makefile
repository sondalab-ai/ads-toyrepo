.PHONY: all test vet

all: vet test

test:
	go test ./...

vet:
	go vet ./...
