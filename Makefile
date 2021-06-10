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
	rm -rf db/badger/data
	docker-compose up postgres -d

start-geth-locally:
	rm -rf e2e/geth-data/geth
	geth --datadir e2e/geth-data --dev --dev.period 1 --http --ws

setup-geth:
	rm -rf e2e/geth-data/geth
	docker-compose up ethereum-node

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
	mkdir -p "e2e-data"
	HUBBLE_E2E=docker go test -v -tags e2e ./e2e
	rm -r "e2e-data"

test-commander-locally: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run TestCommander ./e2e

bench-e2e: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run TestBenchmarkSuite ./e2e

bench-e2e-profile: clean-testcache
	HUBBLE_E2E=in-process go test -cpuprofile cpu.prof -v -tags e2e -run TestBenchmarkSuite ./e2e

.PHONY: 
	install
	clean
	clean-testcache
	compile
	generate
	build
	setup-db
	start-geth-locally
	setup-geth
	update-contracts
	run
	run-prune
	run-dev
	lint
	test
	test-hardhat
	test-e2e
	test-commander-locally
	bench-e2e
	bench-e2e-profile
