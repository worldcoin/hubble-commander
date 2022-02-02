# Hubble Commander

![lines of code](https://img.shields.io/tokei/lines/github/worldcoin/hubble-commander)
[![codecov](https://codecov.io/gh/worldcoin/hubble-commander/branch/main/graph/badge.svg?token=WBPZ9U4TTO)](https://codecov.io/gh/worldcoin/hubble-commander)
[![CI](https://github.com/worldcoin/hubble-commander/actions/workflows/ci.yml/badge.svg)](https://github.com/worldcoin/hubble-commander/actions/workflows/ci.yml)
[![E2E Test](https://github.com/worldcoin/hubble-commander/actions/workflows/e2e.yml/badge.svg)](https://github.com/worldcoin/hubble-commander/actions/workflows/e2e.yml)

## Overview

|                     path                     |                description                |
| -------------------------------------------- | ----------------------------------------- |
| [`commander`](commander)                     | Main application struct                   |
| [`config`](config)                           | Config loading code                       |
| [`contracts`](contracts)                     | Smart contract wrappers                   |
| [`utils`](utils)                             | Utilities                                 |
| [`models`](models)                           | Repository of types                       |
| [`storage`](storage)                         | Storage                                   |
| [`bls`](bls)                                 | BLS Signature library                     |
| [`bls/sdk`](bls/sdk)                         | BLS Wrapper for Mobile client             |
| [`db`](db)                                   | Database abstraction                      |
| [`api`](api)                                 | API Package                               |
| [`eth`](eth)                                 | Ethereum client                           |
| [`main`](main)                               | Command line interface                    |
| [`e2e`](e2e)                                 | End-to-end tests                          |

<!-- Above table extracted using
#!/usr/bin/env python3
from pathlib import Path
from pytablewriter import MarkdownTableWriter
import re

value_matrix = []
for path_object in Path(".").glob('**/Readme.md'):
    if 'hubble-contracts' in f"{path_object}":
        continue
    print(f"Parsing {path_object}")
    dir = f"{path_object.parent}"
    readme = path_object.read_text()
    title = re.match("#\s*(.*?)\n", readme).group(1)
    value_matrix += [[f"[`{dir}`]({dir})", f"{title}"]]

MarkdownTableWriter(
    headers=["path", "description"],
    value_matrix=value_matrix,
    margin=1
).write_table()
-->

## Prerequisites

### Contract bindings

In order to generate Go bindings for smart contracts `abigen` tool needs to be installed locally. 
It comes along with Geth which can be installed on macOS using:

```shell
brew tap ethereum/ethereum
brew install ethereum
```

For other environments refer to: <https://geth.ethereum.org/docs/install-and-build/installing-geth>

You also need python3 installed: <https://www.python.org/>

### golangci-lint

For the lint script to work `golangci-lint` must be installed locally. On macOS run:

```shell
brew install golangci-lint
brew upgrade golangci-lint
```

### mdBook documentation

The documentation server requires Rust and mdBook to be installed.

To install Rust and Cargo on MacOS, type the following:

```shell
brew install rust
```

Then install mdBook and mdbook-toc preprocessor with Cargo:

```shell
cargo install mdbook mdbook-toc
```

Afterwards, run the server:

```shell
make run-docs
```

For other environments refer to: <https://golangci-lint.run/usage/install/#local-installation>

## Running the commander

The commander can be started by running the binary with a `start` subcommand, e.g. `commander start`, and it requires deployed smart
contracts to work. It can connect to said smart contracts by fetching their addresses either from a chain spec file or from an already
running commander. The path to a chain spec file and the url of a remote commander node can be set in the config file (
see `commander-config.example.yaml` file for reference)
or with env variables:

```shell
# Environment variables
HUBBLE_BOOTSTRAP_CHAIN_SPEC_PATH=chain-spec.yaml
HUBBLE_BOOTSTRAP_NODE_URL=http://localhost:8080
```

The smart contracts can be deployed by using the binary with a `deploy` subcommand, e.g. `commander deploy`.
The subcommand uses its own config (see `deployer-config.example.yaml` file for reference).
After a successful deployment, a chain spec file will be generated which can be used to start the commander.
Additionally, the path to a chain spec file can be provided with `file` flag, e.g. `commander deploy -file chain-spec.yaml`.

## Scripts

There is a number of scripts defined in the Makefile:

* `make install` - install Go dependencies
* `make clean` - remove build artifacts
* `make clean-testcache` - remove cached test results 
* `make compile` - build artifacts
* `make generate` - generate bindings for smart contracts
* `make build` - clean and build artifacts
* `make start-geth-locally` - start a new instance of Go-Ethereum node
* `make setup-geth` - create and run a Docker container with Go-Ethereum node
* `make update-contracts` - update the `hubble-contracts` git submodule
* `make deploy` - deploys the smart contracts and generates `chain-spec.yaml` file required for running the commander
* `make run` - run the compiled binary with `start` flag
* `make run-dev` - run-prune without transaction signature validation
* `make start-dev` - deploy and run-dev
* `make export-state` - exports state leaves to a file
* `make export-accounts` - exports accounts to a file
* `make lint` - run linter
* `make test` - run all tests unit tests
* `make run-docs` - render and preview docs by serving it via HTTP
* `make clean-docs` - delete the generated docs and any other build artifacts

Other scripts defined in the Makefile file are used on the CI.

## Running commander with Go-Ethereum (Geth)

Start local Geth either natively or with Docker:

```shell
# Starts native Geth
make start-geth-locally

# Starts Geth in Docker container
make setup-geth
```

Use the following config to make commander scripts connect to the local node:

```shell
HUBBLE_ETHEREUM_RPC_URL=ws://localhost:8546
HUBBLE_ETHEREUM_CHAIN_ID=1337
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
```

Deploy smart contracts:
```shell
make deploy
```

Start commander:
```shell
make run
```

## Running commander image on Docker for Mac

Create `.env.docker` file and set necessary env variables:

```shell
HUBBLE_ETHEREUM_RPC_URL=ws://docker.for.mac.localhost:8546
HUBBLE_ETHEREUM_CHAIN_ID=1337
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
```

Create `chain-spec` directory with:

```shell
mkdir chain-spec
```

Then run this command to deploy the smart contracts and create a chain spec file:

```shell
docker run -it -v $(pwd)/chain-spec:/go/src/app/chain-spec -p 8080:8080 --env-file .env.docker ghcr.io/worldcoin/hubble-commander:latest deploy -file /go/src/app/chain-spec/chain-spec.yaml
```

Afterwards, run this to start the commander:

```shell
docker run -it -v $(pwd)/chain-spec:/go/src/app/chain-spec -p 8080:8080 --env-file .env.docker ghcr.io/worldcoin/hubble-commander:latest start
```

## Running E2E tests and benchmarks
Start Geth in a separate terminal window:
```shell
make start-geth-locally
# OR
make setup-geth
```
Run E2E tests:
```shell
go test -v -tags e2e ./e2e
```

Run E2E benchmarks:
```shell
go test -v -tags e2e ./e2e/bench -timeout 1200s
```

## Running locally

**Step 1.** Build the Hubble commander container image:

```
docker build --tag hubble .
```

**Step 2.** Start a fresh Geth testnet instance (keep running in background).

```
docker run --rm -t -p 8545-8546:8545-8546 ethereum/client-go:stable \
    --dev --dev.period=1 --http --http.addr=0.0.0.0 --ws --ws.addr=0.0.0.0
```

**Step 3.** Fund the Hubble deployer/operator Ethereum account.

```
docker run --rm -t ethereum/client-go:stable \
    attach ws://docker.for.mac.localhost:8546 \
    --exec 'eth.sendTransaction({ 
        from:  eth.accounts[0],
        to:    "0xCd2a3d9f938e13Cd947eC05ABC7fe734df8DD826", 
        value: web3.toWei(1000, "ether") 
    })'
```

**Step 4.** Deploy contracts on Ethereum node using `genesis.yaml` and generate `chain-spec.yaml`.

```
touch chain-spec.yaml
docker run --rm -ti \
    -e HUBBLE_LOG_LEVEL=debug -e HUBBLE_LOG_FORMAT=text \
    -e HUBBLE_ETHEREUM_RPC_URL=ws://docker.for.mac.localhost:8546 \
    -e HUBBLE_ETHEREUM_CHAIN_ID=1337 \
    -e HUBBLE_ETHEREUM_PRIVATE_KEY=c85ef7d79691fe79573b1a7064c19c1a9819ebdbd1faaab1a8ec92344438aaf4 \
    -v $(pwd)/genesis.yaml:/genesis.yaml:ro \
    -v $(pwd)/chain-spec.yaml:/chain-spec.yaml:rw \
    hubble deploy
```

**Step 5.** Run Hubble commander on Ethereum node with `chain-spec.yaml`

```
docker run --rm -ti -p 8080:8080 \
    -e HUBBLE_LOG_LEVEL=debug -e HUBBLE_LOG_FORMAT=text \
    -e HUBBLE_ETHEREUM_RPC_URL=ws://docker.for.mac.localhost:8546 \
    -e HUBBLE_ETHEREUM_CHAIN_ID=1337 \
    -e HUBBLE_ETHEREUM_PRIVATE_KEY=c85ef7d79691fe79573b1a7064c19c1a9819ebdbd1faaab1a8ec92344438aaf4 \
    -e HUBBLE_API_AUTHENTICATION_KEY=89ca4560ec5925f271359196972d762d \
    -e HUBBLE_BOOTSTRAP_CHAIN_SPEC_PATH=/chain-spec.yaml \
    -v $(pwd)/chain-spec.yaml:/chain-spec.yaml:ro \
    hubble
```

## Profiling

To profile batch creation run this command:

```shell
HUBBLE_E2E=in-process go test -cpuprofile creation.prof -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchMixedCommander ./e2e/bench
```

To profile batch syncing run this command:

```shell
HUBBLE_E2E=in-process go test -cpuprofile sync.prof -v -tags e2e -run BenchmarkTransactionsSuite/TestBenchSyncCommander ./e2e/bench
```
