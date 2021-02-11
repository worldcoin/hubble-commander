# Hubble Commander

## Prerequisites

In order to run lint script `golangci-lint` must be installed locally. On macOS run:

```bash
brew install golangci-lint
brew upgrade golangci-lint
```

For other environments refer to: https://golangci-lint.run/usage/install/#local-installation

## Install PostgreSQL

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
