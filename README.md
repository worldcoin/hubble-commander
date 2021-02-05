# Hubble Commander

## Prerequisites

In order to run lint script `golangci-lint` must be installed locally. On macOS run:

```bash
brew install golangci-lint
brew upgrade golangci-lint
```

For other environments refer to: https://golangci-lint.run/usage/install/#local-installation

## Scripts

There are a couple of scripts defined in the Makefile:

* `make install` - install Go dependencies
* `make clean` - remove build artifacts
* `make compile` - build artifacts
* `make build` - clean and build artifacts
* `make run` - run the compiled binary
* `make lint` - run linter
* `make test` - run all tests
