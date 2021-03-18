# Hubble Commander

## Prerequisites

### Bindings
In order to generate Go bindings for smart contracts `abigen` tool needs to be installed locally. 
It comes along with Geth which can be installed on macOS using:
```bash
brew tap ethereum/ethereum
brew install ethereum
```
For other environments refer to: https://geth.ethereum.org/docs/install-and-build/installing-geth

You also need python3 installed: https://www.python.org/

### golangci-lint

For the lint script to work `golangci-lint` must be installed locally.
On macOS run:
```bash
brew install golangci-lint
brew upgrade golangci-lint
```
For other environments refer to: https://golangci-lint.run/usage/install/#local-installation

### PostgreSQL

You can either install the PostgreSQL locally or use Docker for that:
```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_USER=hubble -e POSTGRES_PASSWORD=root -d postgres
```

## Scripts

There are a couple of scripts defined in the Makefile:

* `make install` - install Go dependencies
* `make clean` - remove build artifacts
* `make compile` - build artifacts
* `make generate` - generate bindings for smart contracts
* `make build` - clean and build artifacts
* `make setup-db` - create and run a Docker container with postgres
* `make stop-db` - stop the postgres container
* `make start-db` - start the postgres container
* `make teardown-db` - stop and remove the postgres container
* `make run` - run the compiled binary
* `make lint` - run linter
* `make test` - run all tests

## Running with Ganache

Start ganache cli in a separate terminal. Use the following config to make the node connect to the local instance:

```shell
ETHEREUM_RPC_URL=ws://127.0.0.1:8545
ETHEREUM_CHAIN_ID=1616067554748
ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
```
