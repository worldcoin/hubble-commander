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
	rm -rf db/badger/data

update-contracts:
	git submodule update --remote

run:
	go run ./main/main.go

run-prune:
	go run ./main/main.go -prune

run-dev:
	go run ./main/main.go -prune -dev

start-geth-locally:
	rm -rf e2e/geth-data/geth
	geth --datadir e2e/geth-data --dev --dev.period 1 --http --ws

setup-geth:
	docker run --name ethereum-node -d -v $(CURDIR)/e2e/geth-data:/root/ethereum \
				-p 8545:8545 -p 8546:8546 -p 30303:30303 \
				ethereum/client-go --datadir /root/ethereum \
				--dev --dev.period 1 --http --http.addr 0.0.0.0 --ws --ws.addr 0.0.0.0

stop-geth:
	docker stop ethereum-node

start-geth:
	docker start ethereum-node

teardown-geth:
	docker rm ethereum-node
	rm -rf e2e/geth-data/geth

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
	stop-db
	start-db
	teardown-db
	update-contracts
	run
	run-prune
	run-dev
	start-geth-locally
	setup-geth
	stop-geth
	start-geth
	teardown-geth
	lint
	test
	test-hardhat
	test-e2e
	test-e2e-locally
