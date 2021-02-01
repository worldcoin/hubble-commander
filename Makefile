install:
	go get -v -t -d ./...

clean:
	rm -rf build

build: clean
	mkdir -p build
	go build -o build/hubble ./cmd

lint:
	golangci-lint run ./...

test:
	go test -v ./...

.PHONY: install clean build lint test
