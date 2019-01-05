all: test build
ci: pre-build test

pre-build:
	go get .

build:
	go build .

test: 
	go vet .
	go test ./... -cover
