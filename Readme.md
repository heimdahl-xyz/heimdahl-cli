# Heim-cli:  Realtime Event Listener & Indexer CLI

**EVM Event Listener & Indexer**
is a command-line interface (CLI) tool written in Go.

## Features

- **Create Event Listeners by Contract Address**: Register Event Listener by EVM contract address or contract address and ABI.
- **Listen to Realtime Events**: Listen to events via Websocket API.
- **Query Indexed Blockchain Events via REST API**: Retrieve Event data already indexed on backend.
- **Multiple Blockchain Support**: Supports Ethereum, Arbitrum, and more(in development).
- **Event Filtering**: Filter events by block number, transaction hash, event type, or contract address.
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
go build -o bin/heim-cli main.go
or
make build 
```

## Run CLI
```bash
➜  heimdahl-cli git:(master) ✗ bin/heim-cl --help
A CLI client for interacting with the Heimdahl event listener API

Usage:
  heim-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create a new event listener
  get         Get an event listener by address
  help        Help about any command
  list        List all event listeners
  listen      Listen to a WebSocket connection

Flags:
  -h, --help            help for heim-cli
  -H, --host string     Host URL for the API server (default "https://api.heimdahl.xyz")
  -W, --wsHost string   WSHost URL for the API server (default "wss://api.heimdahl.xyz")

Use "heim-cli [command] --help" for more information about a command.
```

### Create event listener
```bash
 bin/heim-cli create --address=0x15b7c0c907e4C6b9AdaAaabC300C08991D6CEA05 --name="GEL Token" --network="mainnet" --chain="ethereum"
Successfully created event listener for contract 0x15b7c0c907e4C6b9AdaAaabC300C08991D6CEA05%
```

### List event listeners
```bash
 bin/heim-cli list
➜  heimdahl-cli git:(master) ✗ bin/heim-cli list

Network: mainnet
Contract Name: USDT Token
Contract Address: 0xdAC17F958D2ee523a2206206994597C13D831ec7

Network: mainnet
Contract Name: DyDx
Contract Address: 0x92D6C1e31e14520e676a687F0a93788B716BEff5

Network: mainnet
Contract Name: AAVE
Contract Address: 0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2
```

