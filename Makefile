all: test build run-manual
ci: pre-build test

pre-build:
	go get .

build:
	go build .

test: 
	go tool vet .
	go test ./... -cover

run-manual:
	./rdoc

