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

config: config.template.yaml
	cp config.template.yaml config.yaml

.PHONY: install clean compile build run lint test
