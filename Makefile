all: fmt test

fmt:
	gofmt -w ./webkit2/..

test: fmt
	go test ./webkit2

clean:
	rm -rf vendor bin
