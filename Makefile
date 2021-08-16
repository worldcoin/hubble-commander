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
	rm -rf db/badger/data/hubble
	docker-compose up -d postgres

start-geth-locally:
	rm -rf e2e/geth-data/geth
	geth --datadir e2e/geth-data --dev --dev.period 1 --http --ws

setup-geth:
	rm -rf e2e/geth-data/geth
	docker-compose up geth

update-contracts:
	git submodule update --remote

run:
	go run ./main/main.go

run-prune:
	HUBBLE_BOOTSTRAP_PRUNE=true go run ./main/main.go start

run-dev:
	HUBBLE_BOOTSTRAP_PRUNE=true HUBBLE_ROLLUP_DISABLE_SIGNATURES=true go run ./main/main.go start

lint:
	golangci-lint run ./...

test:
	go test -p 1 -v ./...

test-hardhat:
	go test -v -tags hardhat -run TestWalletHardhatTestSuite ./bls

test-e2e: clean-testcache
	mkdir -p "e2e-data"
	go test -v -tags e2e ./e2e
	rm -r "e2e-data"

test-e2e-in-process: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e ./e2e

test-commander-locally: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run TestCommander ./e2e

bench-e2e: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run TestBenchmarkSuite ./e2e

bench-creation-profile: clean-testcache
	HUBBLE_E2E=in-process go test -cpuprofile creation.prof -v -tags e2e -run TestBenchmarkSuite/TestBenchCommander ./e2e

bench-sync-profile: clean-testcache
	HUBBLE_E2E=in-process go test -cpuprofile sync.prof -v -tags e2e -run TestBenchmarkSuite/TestBenchSyncCommander ./e2e

measure-dispute-gas: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e -run TestMeasureDisputeGasUsage ./e2e

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
	test-e2e-in-process
	test-commander-locally
	bench-e2e
	bench-creation-profile
	bench-sync-profile
	measure-dispute-gas
