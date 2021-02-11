install:
	go get -v -t -d ./...

clean:
	rm -rf build

compile:
	mkdir -p build
	go build -o build/hubble ./main

generate:
	cd hubble-contracts && npm install
	cd hubble-contracts && npm run compile
	go generate

build: clean compile

run:
	./build/hubble

lint:
	golangci-lint run ./...

test:
	go test -v ./...

config:
	cp config.template.yaml config.yaml

.PHONY: install clean compile generate build run lint test config
