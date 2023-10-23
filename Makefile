.PHONY: test cov lint clean
all: test lint

test:
	go test -race -cover -timeout=5s ./...

cov:
	@go test -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@go tool cover -func coverage.out

clean:
	rm -f coverage.out
	go clean

lint:
	golangci-lint run
