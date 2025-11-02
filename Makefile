run:
	@go run .

build:
	@go build -o bin/relay

test:
	@go test ./...

dev:
	air 
