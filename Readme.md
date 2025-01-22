# Heimdahl cli: Blockchain data aggregator command line 

![Heimdahl-logo](https://github.com/heimdahl-xyz/heimdahl-cli/blob/572e8557e1dc6181443db5a21123a295b365fb51/static/heimdahl-logo.png?raw=true)


**Heimdahl CLI**
is a general purpose command line tool that aims to provide convenient way to access blockchain data.

## Features

- **Query Indexed Blockchain Events via REST API**: Retrieve smart contracts events data already pre-indexed on backend.
- **Direct sourcing from in-house Ethereum nodes**: Get data directly from blockchain without any middlemen (Alchemy,
  Infura etc) with more nodes to come.
- **Unpacked data** : Get decoded data from events with parsing via ABI.
- **Unified API**: Access multiple chains events data with a single API.
- **Zero Configuration**: No need to run your own infrastructure or indexers.
- **REST API**: Access blockchain data via REST API to eliminate need for complex GraphQL querying
- **Built in Go**: Fast, portable, and efficient CLI with zero configuration and underlying dependency.
- **Simple Commands**: Get started with simple commands powered by [Cobra](https://github.com/spf13/cobra).

### Features in development

- **Fully in-house Solana node**: Access Solana data with same ease.
- **Listen to Realtime Events**: Listen to events via Websocket API.
- **Multiple Blockchain Support**: Supports for Arbitrum and more(in development).
- **Advanced filtering** Filter events by block number, timestamp, and more.
- **Block, Transactions and Receipts Indexing** Indexing of blocks, transactions and receipts for advanced analytics.
- **Realtime Websocket API**: Listen to events in real-time via Websocket API.

- and many more to come ;)

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

or you can ...

## Build from source

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

### List events

#### List Uniswap V3 Mint events on Ethereum mainnet

```bash
$ heimdahl event list 0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640 Mint # Uniswap 
V3 Mint on ethereum
           
BLOCK#     | TIMESTAMP       | TRANSACTION_HASH                                                  | EVENT_DATA                                          
----------------------------------------------------------------------------------------------------                                                          21446044 | 2025-01-17T18:54:55Z | 0xb88af7b04df4d1d7a7fd777bcc31060ad0346a0da3706a32d15d4ddcf296e334 | amount: 393581927855607704795, amount1: 361788791691837
3918, tickLower: 194800, amount0: 11576768988867, owner: 0xA69babEF1cA67A37Ffaf7a485DfFF3382056e78C, sender: 1.4024810831325053e48, tickUpper: 194810         21446018 | 2025-01-17T18:54:55Z | 0x91fb7b4208e216dda90c6f83d042b528d4208a86e7bddf8e1b99e473d19f6c1e | amount0: 145884955, tickLower: 193890, tickUpper: 19808
0, amount: 16340383564036, owner: 0xC36442b4a4522E871399CD717aBDD847Ab11FE88, sender: 1.1154890857106194e48, amount1: 12253693449477527                       21445992 | 2025-01-17T18:54:55Z | 0x72de3607664cab8ab4179b4b1f581cfd8b1b219944b705e9ff7ce8248ae545a4 | sender: 1.1154890857106194e48, tickUpper: 198080, amoun
t: 51397277183042, amount0: 458518761, owner: 0xC36442b4a4522E871399CD717aBDD847Ab11FE88, tickLower: 193890, amount1: 38643713455569501                       21445953 | 2025-01-17T18:54:55Z | 0x9d12b2d2b68c7bc0a920d1d950639f11f940af43ab711d5266dc7c3cf6efe0d7 | amount0: 601856883798, amount1: 0, amount: 204088654544
29306880, owner: 0x1f2F10D1C40777AE1Da742455c65828FF36Df387, sender: 1.7802831409602703e47, tickLower: 194770, tickUpper: 194780               
...
```

#### List WETH Transfers on Ethereum mainnet

```bash
$ heimdahl event list 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2 Transfer 

BLOCK#     | TIMESTAMP       | TRANSACTION_HASH                                                  | EVENT_DATA                                                 ----------------------------------------------------------------------------------------------------                                                          
21679357 | 2025-01-22T10:13:38Z | 0x2b1cbb5891849ae99a1436451add7c2212b0b782c8b95259d52b515f20afef9f | dst: 0xB86E490E72F050c424383d514362Dc61DaBB1Cc3, src: 0xdad17D7E3Abbebe1ea5782962398113422F10EE0, wad: 300852659872045380                                                                                            
21679357 | 2025-01-22T10:13:38Z | 0x318a0434558e1c6b70c3b60eb58f74751ac8559b7dca4aa3c351166e75f57bff | dst: 0x95190AaF90dd499E87068C68b90352526993c1A7, src: 0
x7a250d5630B4cF539739dF2C5dAcb4c659F2488D, wad: 200000000000000000
21679357 | 2025-01-22T10:13:38Z | 0xe3157a424eb975be6d25376f2021cdbcc0d3ec370d2c8d93a81fa12928e1f8c3 | dst: 0x000000fee13a103A10D593b9AE06b3e05F2E7E1c, src: 0
x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD, wad: 73468283229723
21679357 | 2025-01-22T10:13:38Z | 0xcc60fbf73b6b999efa3efabcae87a18a3b1499c2f70a4dcf6dc3ad7db9d252bf | dst: 0x468795E031c173942C9387AEd0a302E26bDD0460, src: 0
x7D0CcAa3Fac1e5A943c5168b6CEd828691b46B36, wad: 10074112441007251
21679357 | 2025-01-22T10:13:38Z | 0xa161fcaac3f5d1b90f8b95fdf1e86baac469fe9fad877690e668efc7a858b6e4 | dst: 0x468795E031c173942C9387AEd0a302E26bDD0460, src: 0
x3b4c1914CCbb2ADb030a58A2cC378b7016C9d5E8, wad: 53024641812615367

```

### How to start?

We're actively working on enabling users to obtain their API keys independently. In the meantime, you can gain early
access by joining our [Discord Server](https://discord.gg/kZKA977B).

You can also follow us on [LinkedIn](https://www.linkedin.com/company/heimdahl-xyz/?viewAsMember=true) to stay updated.
