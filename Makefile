build:
	@go build -o ./bin/main ./cmd/web

run: build
	@./bin/main $(args)

test:
	@go test -v ./...
