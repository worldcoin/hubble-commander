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

setup-db:
	docker run --name hubble-postgres -p 5432:5432 -e POSTGRES_USER=hubble -e POSTGRES_PASSWORD=root -d postgres

stop-db:
	docker stop hubble-postgres

start-db:
	docker start hubble-postgres

teardown-db: stop-db
	docker rm hubble-postgres

run:
	./build/hubble

lint:
	golangci-lint run ./...

test:
	go test -v ./...

.PHONY: install clean compile generate build setup-db stop-db start-db teardown-db run lint test
