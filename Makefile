install:
	go get -v -t -d ./...

clean:
	rm -rf build

build: clean compile

compile:
	mkdir -p build
	go build -o build/hubble ./main

run:
	./build/hubble

lint:
	golangci-lint run ./...

test:
	go test -v ./...

.PHONY: install clean build compile run lint test
