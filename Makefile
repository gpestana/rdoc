all: test build run-manual

build:
	go build .

test: 
	go test .

run-manual:
	./crdt-json

