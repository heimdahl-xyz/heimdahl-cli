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

### Stream Transfer events for USDT in Ethereum 
```
➜  heimdahl-cli git:(master) ✗ bin/heimdahl stream listen   --address=0xdAC17F958D2ee523a2206206994597C13D831ec7 --event=Transfer

CHAIN      | NETWORK    | BLOCK#   | BLOCK_HASH      | TIMESTAMP           | CONTRACT        | TRANSACTION_HASH | EVENT_DATA                                  
------------------------------------------------------------------------------------------------------------------------                      ethereum   | mainnet    | 21172635 | 0x1f934caba9a749257dd2d3492087ac03d8544a64c115159011353a09e9f0811e | 2024-11-12 15:56:13 | 0xdAC17F958D2ee523a22062069945
97C13D831ec7 | to: 0x6c207dCD2c7bD0b1f6cC5Cb3319D54662a1a62a4, transactionIndex: 7, value: 2.051946238e+09, from: 0xF955C57f9EA9Dc8781965FEaE0b6A2acE2BAD6f3
ethereum   | mainnet    | 21172635 | 0x1f934caba9a749257dd2d3492087ac03d8544a64c115159011353a09e9f0811e | 2024-11-12 15:56:13 | 0xdAC17F958D2ee523a22062069945
97C13D831ec7 | from: 0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852, transactionIndex: 10, to: 0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD, value: 1.345181186e+09 
ethereum   | mainnet    | 21172635 | 0x1f934caba9a749257dd2d3492087ac03d8544a64c115159011353a09e9f0811e | 2024-11-12 15:56:13 | 0xdAC17F958D2ee523a22062069945
97C13D831ec7 | to: 0x480f0D4887eD4F16D2299031dffEc90782826269, from: 0xF955C57f9EA9Dc8781965FEaE0b6A2acE2BAD6f3, value: 2.4994062e+07, transactionIndex: 7  
```


### Replay Ethereum USDT Approvals 
```
heimdahl stream replay   --address=0xdAC17F958D2ee523a2206206994597C13D831ec7 --
event=Approval                                                                                                                                                BLOCK#     | BLOCK_HASH                                                        | TIMESTAMP | CONTRACT        | TRANSACTION_HASH    | EVENT_DATA     
----------------------------------------------------------------------------------------------------
21114591 | 0xe9fec20213e8c5c642daf31040400a9dab90ed3f3c980acce6e5330969763fc5 | 2024-11-07T13:22:00Z | 9 | value: 11579208923731619542357098500868790785326998
4665640564039457584007913129639935, spender: 0x881D40237659C251811CEC9c364ef91dC08D300C, owner: 0x03de42d3D23Da88ef3FE72F2569449641BBd49C0    21114591 | 0xe9fec20213e8c5c642daf31040400a9dab90ed3f3c980acce6e5330969763fc5 | 2024-11-07T13:22:00Z | 2 | owner: 0xb6F2D272584052E612Be87F5A5e45a3Cf12b9c1B, 
spender: 0x216B4B4Ba9F3e719726886d34a177484278Bfcae, value: 115792089237316195423570985008687907853269984665640564039457584007913129639935    21114592 | 0x5117bbde3f74e638865f2efea584ae12c22646777962c2d04817f0a610746a82 | 2024-11-07T13:22:01Z | 18 | owner: 0x4a14347083B80E5216cA31350a2D21702aC3650d,
 spender: 0xE592427A0AEce92De3Edee1F18E0157C05861564, value: 0                                                                                21114593 | 0x1ba50e4ef2cd88b23ffe43f64241232c8c370dcc63237a64678c520ff42674b1 | 2024-11-07T13:22:00Z | 13 | spender: 0xC9f93163c99695c6526b799EbcA2207Fdf7D61a
D, owner: 0x1f2F10D1C40777AE1Da742455c65828FF36Df387, value: 121276199680           
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

