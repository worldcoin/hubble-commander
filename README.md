# Hubble Commander

## Prerequisites

### Bindings
In order to generate Go bindings for smart contracts `abigen` tool needs to be installed locally. 
It comes along with Geth which can be installed on macOS using:
```shell
brew tap ethereum/ethereum
brew install ethereum
```
For other environments refer to: https://geth.ethereum.org/docs/install-and-build/installing-geth

You also need python3 installed: https://www.python.org/

### golangci-lint

For the lint script to work `golangci-lint` must be installed locally.
On macOS run:
```shell
brew install golangci-lint
brew upgrade golangci-lint
```
For other environments refer to: https://golangci-lint.run/usage/install/#local-installation

### PostgreSQL

You can either install the PostgreSQL locally or use Docker for that:
```shell
make setup-db
```

## Scripts

There is a number of scripts defined in the Makefile:

* `make install` - install Go dependencies
* `make clean` - remove build artifacts
* `make clean-testcache` - remove cached test results 
* `make compile` - build artifacts
* `make generate` - generate bindings for smart contracts
* `make build` - clean and build artifacts
* `make setup-db` - create and run a Docker container with postgres
* `make stop-db` - stop the postgres container
* `make start-db` - start the postgres container
* `make teardown-db` - stop and remove the postgres container
* `make update-contracts` - update the `hubble-contracts` git submodule
* `make run` - run the compiled binary
* `make run-prune` - clean database and run the compiled binary
* `make run-dev` - run-prune without transaction signature validation
* `make lint` - run linter
* `make test` - run all tests unit tests
* `make test-hardhat` - run all tests with Hardhat dependency
* `make test-e2e` - run E2E tests on a pre-built docker image
* `make test-commander-locally` - run E2E tests against a local commander instance

## Running with Ganache

Start Ganache CLI in a separate terminal:
```shell
npx ganache-cli --account 0xee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82,0x56BC75E2D63100000
```

Use the following config to make commander connect to the local node
```shell
HUBBLE_ETHEREUM_RPC_URL=ws://127.0.0.1:8545
HUBBLE_ETHEREUM_CHAIN_ID=1616067554748
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
```

## Running docker image on Docker for Mac
Create `.env.docker` file and set necessary env variables:
```
HUBBLE_ETHEREUM_RPC_URL=ws://docker.for.mac.localhost:8545
HUBBLE_ETHEREUM_CHAIN_ID=1616067554748
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
HUBBLE_POSTGRES_HOST=docker.for.mac.localhost
HUBBLE_POSTGRES_USER=hubble
HUBBLE_POSTGRES_PASSWORD=root
```

Then run:
```shell
docker run -it --rm -p 8080:8080 --env-file .env.docker ghcr.io/worldcoin/hubble-commander:latest
```

## Running E2E tests against a Docker image

Build the docker image:
```shell
docker build . -t ghcr.io/worldcoin/hubble-commander:latest
```

Export variables from the `.env.docker` to the currently running shell:
```shell
export $(grep -v '#.*' .env.docker | xargs)
```

Run the E2E tests:
```shell
make test-e2e
```
The Docker container will be started and stopped automatically.

## Running E2E tests locally

Start commander in a separate terminal window and wait until it finished bootstrapping:
```shell
make run-prune
```

Run the E2E tests:
```shell
make test-commander-locally
```
