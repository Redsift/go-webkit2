all: fmt test

fmt:
	gofmt -w ./webkit2/..

test: fmt
	CC=gcc-4.9 go test ./webkit2

clean:
	rm -rf vendor bin
