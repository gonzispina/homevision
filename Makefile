# cross parameters
SHELL:=/bin/bash -O extglob
BINARY=bin/homevision
VERSION=0.1.0

LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Build step, generates the binary.
run:
	@go build ${LDFLAGS} -o ${BINARY} cmd/main.go
	@./bin/homevision