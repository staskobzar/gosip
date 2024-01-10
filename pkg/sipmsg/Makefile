all: ragel

ragel:
	ragel -Z -G2 -o parse_uri.go parse_uri.rl

test:
	ragel -Z -T0 -o parse_uri.go parse_uri.rl
	go test -race -cover -timeout=5s

bench:
	ragel -Z -G2 -o parse_uri.go parse_uri.rl
	go test -benchmem -bench=.

cov:
	ragel -Z -T0 -o parse_uri.go parse_uri.rl
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
