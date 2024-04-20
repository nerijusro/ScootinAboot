build:
	@go build -o bin/scootinAboot cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/scootinAboot