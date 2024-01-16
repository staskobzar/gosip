all: sipmsg/ragel

sipmsg/ragel:
	ragel -Z -G2 -o parse_uri.go parse_uri.rl
	ragel -Z -G2 -o parse.go parse.rl

test:
	ragel -Z -T0 -o parse_uri.go parse_uri.rl
	ragel -Z -T0 -o parse.go parse.rl
	go test -race -cover -timeout=5s ${TESTOPTS}

bench:
	ragel -Z -G2 -o parse_uri.go parse_uri.rl
	ragel -Z -G2 -o parse.go parse.rl
	go test -benchmem -bench=.

cov:
	ragel -Z -T0 -o parse_uri.go parse_uri.rl
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
