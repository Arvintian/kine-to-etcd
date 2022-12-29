
build: k2e kinectl

k2e: $(shell find cmd/k2e -type f -name '*.go')
	go build -o bin/k2e cmd/k2e/*.go

kinectl: $(shell find cmd/kinectl -type f -name '*.go')
	go build -o bin/kinectl cmd/kinectl/*.go