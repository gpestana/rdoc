tests:
	go test ./... -v -cover

lint:
	golangci-lint run -E gofmt -E golint --exclude-use-default=false
