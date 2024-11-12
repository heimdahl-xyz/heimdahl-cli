# Heimdahl cli: Command line tool for Heimdahl API

**EVM events indexer and listener**
is a command-line interface (CLI) tool.

## Features
- **Listen to Realtime Events**: Listen to events via Websocket API.
- **Query Indexed Blockchain Events via REST API**: Retrieve Event data already indexed on backend.
- **Multiple Blockchain Support**: Supports Ethereum, Arbitrum, and more(in development).
- **Built in Go**: Fast, portable, and efficient CLI.

## Installation

### Prerequisites

- **Go**: Ensure Go is installed (version 1.21+).
- **Make**: GNU Make 3.81(optional)


### Clone the Repository

```bash
git clone git@github.com:heimdahl-xyz/heimdahl-cli.git
cd heimdahl-cli
```
## Build

```bash
go build -o bin/heimdahl main.go
or
make build 
```

## Run CLI
```bash
➜  heimdahl-cli git:(master) ✗ bin/heimdahl --help
A CLI client for interacting with the Heimdahl event listener API

Usage:
  heimdahl [command]

Available Commands:
  chain       Chain subcommands
  completion  Generate the autocompletion script for the specified shell
  contract    Contract subcommands
  help        Help about any command
  stream      Stream subcommands

Flags:
  -K, --apiKey string   API Key for connection to server (default "test1")
  -h, --help            help for heimdahl
  -H, --host string     Host URL for the API server (default "api.heimdahl.xyz")
      --secure          Use secure connection to server (default true)

Use "heimdahl [command] --help" for more information about a command.
```

### List event listeners
```bash
➜  heimdahl-cli git:(master) ✗ bin/heimdahl contract list
               
NETWORK    CONTRACT NAME   CONTRACT ADDRESS                                                                                                                   
-------------------------------------------------------------------------------                      
arbitrum   | mainnet    | Radiant Token   | 0x0C4681e6C0235179ec3D4F4fc4DF3d14FDD96017
arbitrum   | mainnet    | Radiant Token   | 0x3082CC23568eA640225c2467653dB90e9250AaA0
arbitrum   | mainnet    | SushiSwap V2 Factory | 0xc35DADB65012eC5796536bD9864eD8773aBc74C4
arbitrum   | mainnet    | USDC Coin       | 0xaf88d065e77c8cC2239327C5EDb3A432268e5831
ethereum   | mainnet    | OKX DEX         | 0x40aA958dd87FC8305b97f2BA922CDdCa374bcD7f
ethereum   | mainnet    | DYDX            | 0x92D6C1e31e14520e676a687F0a93788B716BEff5
ethereum   | mainnet    | Ethena USD      | 0x4c9edd5852cd905f086c759e8383e09bff1e68b3
ethereum   | mainnet    | Band Protocol (ORACLE) | 0xBA11D00c5f74255f56a5E366F4F77f5A186d7f55
ethereum   | mainnet    | API 3 Protocol (ORACLE) | 0xa0AD79D995DdeeB18a14eAef56A549A04e3Aa1Bd
arbitrum   | mainnet    | API 3 Protocol (ORACLE) | 0xb015ACeEdD478fc497A798Ab45fcED8BdEd08924
ethereum   | mainnet    | 1inch aggregator | 0x1111111254EEB25477B68fb85Ed929f73A960582
ethereum   | mainnet    | 1inch aggregator v2 | 0x07D91f5fb9Bf7798734C3f606dB065549F6893bb
...
```

