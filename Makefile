.PHONY: test cov lint clean
all: test lint

include pkg/sipmsg/package.mk

test:
	go test -race -cover -timeout=5s ./...

cov:
	@go test -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@go tool cover -func coverage.out

prof-block:
	go test -blockprofile block.out ./...
	go tool pprof block.out

prof-mu:
	go test -mutexprofile mutex.out ./...
	go tool pprof mutex.out
	
trace:
	go test -trace trace.out ./...
	go tool trace trace.out
	
clean:
	rm -f coverage.out
	go clean

lint:
	golangci-lint run
