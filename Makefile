install:
	go get -v -t -d ./...

clean:
	rm -rf build

clean-testcache:
	go clean -testcache

compile:
	mkdir -p build
	go build -o build/hubble ./main

generate:
	cd hubble-contracts && npm ci
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

update-contracts:
	git submodule update --remote

run:
	go run ./main/main.go

run-prune:
	go run ./main/main.go -prune

run-dev:
	go run ./main/main.go -prune -dev

lint:
	golangci-lint run ./...

test:
	go test -p 1 -v ./...

test-hardhat:
	go test -v -tags hardhat -run TestWalletHardhatTestSuite ./bls

test-e2e: clean-testcache
	HUBBLE_E2E=docker go test -v -tags e2e ./e2e

test-commander-locally: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run TestCommander ./e2e

bench-e2e: clean-testcache
	go test -v -tags e2e -run TestBenchCommander ./e2e

.PHONY: 
	install
	clean
	clean-testcache
	compile
	generate
	build
	setup-db
	stop-db
	start-db
	teardown-db
	update-contracts
	run
	run-prune
	run-dev
	lint
	test
	test-hardhat
	test-e2e
	test-e2e-locally
