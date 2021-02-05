install:
	go get -v -t -d ./...

clean:
	rm -rf build

compile:
	mkdir -p build
	go build -o build/hubble ./main

build: clean compile

run:
	./build/hubble

lint:
	golangci-lint run ./...

test:
	go test -v ./...

.PHONY: install clean compile build run lint test
