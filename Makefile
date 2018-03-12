all: test build run-manual

build:
	go build .

test: 
	go tool vet .
	go test ./... -cover

run-manual:
	./crdt-json

