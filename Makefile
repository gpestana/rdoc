all: test build run-manual

build:
	go build .

test: 
	go test ./... -cover

run-manual:
	./crdt-json

