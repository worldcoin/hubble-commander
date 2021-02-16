# Hubble Commander

## Prerequisites

### abigen

In order to generate Go bindings for smart contracts `abigen` tool needs to be installed locally. 
It comes along with Geth which can be installed on macOS using:
```bash
brew tap ethereum/ethereum
brew install ethereum
```
For other environments refer to: https://geth.ethereum.org/docs/install-and-build/installing-geth

### golangci-lint

For the lint script to work `golangci-lint` must be installed locally.
On macOS run:
```bash
brew install golangci-lint
brew upgrade golangci-lint
```
For other environments refer to: https://golangci-lint.run/usage/install/#local-installation

### PostgreSQL

You can either install the PostgreSQL locally or use docker for that:
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
* `make run` - run the compiled binary
* `make lint` - run linter
* `make test` - run all tests
