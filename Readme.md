# Heimdahl cli: Command line tool for EVM chain event stream and replay

**EVM events indexer and listener**
is a command-line interface (CLI) tool.

## Features

- **Query Indexed Blockchain Events via REST API**: Retrieve Event data already indexed on backend.
- **Built in Go**: Fast, portable, and efficient CLI.

### Features in development

- **Listen to Realtime Events**: Listen to events via Websocket API.
- **Multiple Blockchain Support**: Supports for Arbitrum and more(in development).

## Installation

You can pick pre-build binaries from the [releases](https://github.com/heimdahl-xyz/heimdahl-cli/releases) page.

```bash
$ wget https://github.com/heimdahl-xyz/heimdahl-cli/releases/download/heimdahl-cli-dc9e278/heimdahl-cli-linux-amd64.tar.gz
$ tar -xzvf heimdahl-linux-amd64.tar.gz 
$ mv heimdahl-linux-amd64 heimdahl
$ ./heimdahl

Heimdahl CLI - Blockchain Data Access Without Infrastructure

Fast access to blockchain events and analytics across multiple chains through simple commands.
Instead of spending months building indexers, start exploring blockchain data in minutes:

Examples:
  heimdahl event list 0xb47e...BBB PunkOffered   # Get specific events
  heimdahl contract show 0x060...6d              # View contract details

Built for developers who need reliable blockchain data without infrastructure overhead.

```

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

## Quickstart

### List available chains

```bash
$ heimdahl chain list   

CHAIN      NETWORK    CHAIN ID
--------------------------------
ethereum   mainnet    1       
ethereum   sepolia    11155111
base       mainnet    8453    
polygon    mainnet    137     
arbitrum   mainnet    42161   
optimism   mainnet    10      
binance    mainnet    56      

```

### List indexed contracts across chains

```
$ heimdahl contract list 

Chain:            ethereum                                                                                                                                    
Network:          mainnet                                                                                                                                     
Contract Name:    Ethena USD                                                   
Contract Address: 0x4c9edd5852cd905f086c759e8383e09bff1e68b3                                                                                                  
Events:                                                                        
  - EIP712DomainChanged                                                                                                                                       
  - MinterUpdated                 
  - OwnershipTransferStarted                                                                                                                                  
  - OwnershipTransferred                                                                                                                                      
  - Transfer                                                                                                                                                  
  - Approval                                                                                                                                                  
--------------------------------------------------------------------------------                                                                              
Chain:            ethereum
Network:          mainnet
Contract Name:    Pudgy Penguins ETH
Contract Address: 0xBd3531dA5CF5857e7CfAA92426877b022e612cf8
Events:
  - OwnershipTransferred
  - Paused
  - Transfer
  - Unpaused
  - Approval
  - ApprovalForAll
  - CreatePenguin
--------------------------------------------------------------------------------
Chain:            ethereum
Network:          mainnet
Contract Name:    USDD
Contract Address: 0x0C10bF8FcB7Bf5412187A595ab97a3609160b5c6
Events:
  - Approval
  - MetaTransactionExecuted
  - RoleAdminChanged
  - RoleGranted
  - RoleRevoked
  - Transfer

```

### List events for a contract

```bash
$ heimdahl event list 
               
...
```
