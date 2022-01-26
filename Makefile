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
	go generate utils/generate.go

build: clean compile

start-geth-locally:
	rm -rf e2e/geth-data/geth
	geth --datadir e2e/geth-data --dev --dev.period 1 --http --ws

setup-geth:
	rm -rf e2e/geth-data/geth
	docker-compose up geth

update-contracts:
	git submodule update --remote

deploy:
	go run ./main deploy

run:
	go run ./main start

run-prune:
	HUBBLE_BOOTSTRAP_PRUNE=true go run ./main start

run-dev:
	HUBBLE_BOOTSTRAP_PRUNE=true HUBBLE_ROLLUP_DISABLE_SIGNATURES=true go run ./main start

start-dev: deploy run-dev

export-state:
	go run ./main export -type=state

export-accounts:
	go run ./main export -type=accounts

lint:
	golangci-lint run --build-tags hardhat,e2e --fix ./...

test:
	go test -v $(CI_FLAGS) ./...

test-hardhat:
	go test -v -tags hardhat $(CI_FLAGS) ./bls/hardhat

test-e2e-in-process: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e $(CI_FLAGS) ./e2e

test-e2e-locally: clean-testcache
	HUBBLE_E2E=local go test -v -tags e2e -run=^$(TEST)$$ ./e2e

bench-e2e-in-process: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e ./e2e/bench -timeout 1200s

bench-e2e-locally: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e -run=^$(TEST)$$ ./e2e/bench

bench-creation-profile: clean-testcache
	HUBBLE_E2E=in-process go test -cpuprofile creation.prof -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchMixedCommander ./e2e/bench

bench-sync-profile: clean-testcache
	HUBBLE_E2E=in-process go test -cpuprofile sync.prof -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchSyncCommander ./e2e/bench

bench-e2e-ci-part-1: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e -run "BenchmarkTransactionsSuite/(?:TestBenchTransfersCommander|TestBenchCreate2TransfersCommander|TestBenchMassMigrationsCommander)" ./e2e/bench

bench-e2e-ci-part-2: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchMixedCommander ./e2e/bench

bench-e2e-ci-part-3: clean-testcache
	HUBBLE_E2E=in-process go test -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchSyncCommander ./e2e/bench

run-docs:
	mdbook serve

clean-docs:
	mdbook clean

.PHONY:
	install
	clean
	clean-testcache
	compile
	generate
	build
	start-geth-locally
	setup-geth
	update-contracts
	deploy
	run
	run-prune
	run-dev
	start-dev
	export-state
	export-accounts
	lint
	test
	test-hardhat
	test-e2e-in-process
	test-e2e-locally
	bench-e2e-in-process
	bench-e2e-locally
	bench-creation-profile
	bench-sync-profile
	bench-e2e-ci-part-1
	bench-e2e-ci-part-2
	bench-e2e-ci-part-3
	run-docs
	clean-docs
