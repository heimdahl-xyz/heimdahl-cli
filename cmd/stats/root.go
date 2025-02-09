package stats

import (
	"fmt"
	"github.com/spf13/cobra"
)

var dummy = `
-------------------------------------------------------------
                    Heimdahl Stats Summary
-------------------------------------------------------------
Chains Indexed:
------------------------------------------------------
| Blockchain Name   | Indexed Blocks | Contracts Indexed | Total Transactions |
------------------------------------------------------
| Ethereum          | 12,500,000     | 15,000            | 2,500,000          |
| Polygon           | 4,500,000      | 8,000             | 1,200,000          |
| Tron              | 3,000,000      | 6,000             | 800,000            |
| Solana            | 1,200,000      | 3,500             | 500,000            |
| Avalanche         | 600,000        | 2,000             | 300,000            |
------------------------------------------------------

Blockchain Block Stats:
------------------------------------------------------
| Blockchain Name   | Total Blocks Indexed | Latest Block # | Last Block Time  |
------------------------------------------------------
| Ethereum          | 12,500,000           | 13,223,851     | 2025-01-16 12:34:56 |
| Polygon           | 4,500,000            | 4,500,500      | 2025-01-16 12:30:22 |
| Tron              | 3,000,000            | 3,000,500      | 2025-01-16 12:28:43 |
------------------------------------------------------

Transaction Stats:
------------------------------------------------------
| Blockchain Name   | Total Transactions Indexed | Avg Tx Value (ETH) | Avg Gas Price (gwei) |
------------------------------------------------------
| Ethereum          | 2,500,000                 | 0.15               | 35                   |
| Polygon           | 1,200,000                 | 0.12               | 25                   |
| Tron              | 800,000                   | 0.10               | 20                   |
------------------------------------------------------

Contract and Event Stats:
-----------------------------------------------------------
| Blockchain Name   | Contracts Indexed | Total Events Indexed |
-----------------------------------------------------------
| Ethereum          | 15,000            | 1,200,000            |
| Polygon           | 8,000             | 500,000              |
| Tron              | 6,000             | 300,000              |
| Solana            | 3,500             | 200,000              |
-----------------------------------------------------------

Contract Event Stats:
---------------------------------------------------------------
| Contract Address    | Event Type   | Total Events Indexed |
---------------------------------------------------------------
| 0x123abc456         | Transfer     | 1,000,000            |
| 0x123abc456         | Mint         | 300,000              |
| 0x987def123         | Burn         | 150,000              |
---------------------------------------------------------------

Remark: The data displayed below is for demonstration purposes only and does not reflect real statistics.
`
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display Heimdahl stats",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dummy)
	},
}
