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
	./utils/fundAccount.sh &
	geth --datadir e2e/geth-data --dev --dev.period 1 --http --ws --http.api "eth,miner" --ws.api "eth,miner"

setup-geth:
	rm -rf e2e/geth-data/geth
	docker-compose up geth

update-contracts:
	git submodule update --remote

deploy:
	go run ./main deploy

run:
	go run ./main start

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
	go test -v ./...

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
	run-dev
	start-dev
	export-state
	export-accounts
	lint
	test
	run-docs
	clean-docs
