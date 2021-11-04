# Hubble Commander

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

### Bindings

In order to generate Go bindings for smart contracts `abigen` tool needs to be installed locally. 
It comes along with Geth which can be installed on macOS using:

```shell
brew tap ethereum/ethereum
brew install ethereum
```

For other environments refer to: <https://geth.ethereum.org/docs/install-and-build/installing-geth>

You also need python3 installed: <https://www.python.org/>

### golangci-lint

For the lint script to work `golangci-lint` must be installed locally.
On macOS run:

```shell
brew install golangci-lint
brew upgrade golangci-lint
```

For other environments refer to: <https://golangci-lint.run/usage/install/#local-installation>

## Running the commander

The commander can be started by running the binary with a `start` subcommand, e.g. `commander start`,
and it requires deployed smart contracts to work. It can connect to said smart contracts by fetching
their addresses either from a chain spec file or from an already running commander. The path to a chain spec file
and the url of a remote commander node can be set in the config file (see `commander-config.example.yaml` file for reference)
or with env variables:

```shell
# Environmental variables
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
* `make run-prune` - clean database and run the compiled binary with `start` flag
* `make run-dev` - run-prune without transaction signature validation
* `make start-dev` - deploy and run-dev
* `make lint` - run linter
* `make test` - run all tests unit tests
* `make test-hardhat` - run all tests that depend on Hardhat node
* `make test-e2e-in-process` - start commander and run all E2E tests in the same process 
* `make test-e2e-locally` - run an E2E test on a local commander instance (set `TEST` env var to the name of the test)
* `make bench-e2e-in-process` - start commander and run all E2E benchmarks in the same process 
* `make bench-e2e-locally` - run an E2E benchmark on a local commander instance (set `TEST` env var to the name of the test)
* `make bench-creation-profile` - start commander and run E2E batch creation benchmark in the same process with CPU profiling
* `make bench-sync-profile` - start commander and run E2E butch sync benchmark in the same process with CPU profiling
* `make bench-e2e-ci-part-1` - start commander and run 1st part (optimized for the CI speed) of E2E benchmarks in the same process
* `make bench-e2e-ci-part-3` - start commander and run 2nd part (optimized for the CI speed) of E2E benchmarks in the same process
* `make bench-e2e-ci-part-4` - start commander and run 3rd part (optimized for the CI speed) of E2E benchmarks in the same process
* `make bench-e2e-ci-part-5` - start commander and run 4th part (optimized for the CI speed) of E2E benchmarks in the same process

## Running with Go-Ethereum (Geth)

Start Geth either locally or with Docker:

```shell
# Starts geth locally
make start-geth-locally

# Stars geth in docker container
make setup-geth
```

Use the following config to make commander connect to the local node:

```shell
HUBBLE_ETHEREUM_RPC_URL=ws://127.0.0.1:8546
HUBBLE_ETHEREUM_CHAIN_ID=1337
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
```

## Running docker image on Docker for Mac

Create `.env.docker` file and set necessary env variables:

```shell
HUBBLE_ETHEREUM_RPC_URL=ws://docker.for.mac.localhost:8546
HUBBLE_ETHEREUM_CHAIN_ID=1337
HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
HUBBLE_ROLLUP_MIN_TXS_PER_COMMITMENT=32
HUBBLE_ROLLUP_MAX_TXS_PER_COMMITMENT=32
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
make test-e2e-locally
```
