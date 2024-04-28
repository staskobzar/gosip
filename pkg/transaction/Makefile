BENCHOPTS ?= .

all: test

test:
	go test -race -cover -timeout=5s ${TESTOPTS}

bench:
	go test -benchmem -bench=$(BENCHOPTS)

cov:
	go test -coverprofile=coverage.out 
	go tool cover -html=coverage.out
