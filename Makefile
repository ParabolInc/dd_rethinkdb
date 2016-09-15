REV = $(shell git rev-parse --short HEAD)

build:
	GOOS=linux go build -ldflags "-X version.GitCommit=$(REV)" -o bin/dd_rethinkdb .
	@zip dd_rethinkdb.zip bin/dd_rethinkdb
	@rm -rf bin
