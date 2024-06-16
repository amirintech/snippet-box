build:
	@go build -o ./bin/main ./cmd/web

run: build
	@./bin/main

test:
	@go test -v ./...
