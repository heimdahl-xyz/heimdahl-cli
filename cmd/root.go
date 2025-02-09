package cmd

import (
	"fmt"
	chain "github.com/heimdahl-xyz/heimdahl-cli/cmd/chains"
	"github.com/heimdahl-xyz/heimdahl-cli/cmd/contract"
	"github.com/heimdahl-xyz/heimdahl-cli/cmd/event"
	"github.com/heimdahl-xyz/heimdahl-cli/cmd/stats"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/heimdahl-xyz/heimdahl-cli/cmd/fungibles"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "heimdahl",
	Short: "heimdahl - Blockchain Data Access CLI\n",
	Long: `
Heimdahl CLI - Blockchain Data Access Without Infrastructure

Fast access to blockchain events and analytics across multiple chains through simple commands.
Instead of spending months building indexers, start exploring blockchain data in minutes:

Examples:
  heimdahl event list 0xb47e...BBB PunkOffered   # Get specific events
  heimdahl contract show 0x060...6d              # View contract details

Built for developers who need reliable blockchain data without infrastructure overhead.
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add a global flag for host
	RootCmd.PersistentFlags().StringVarP(&config.Config.APIURL, "host", "H", "api.heimdahl.xyz", "Host URL for the API server")
	RootCmd.PersistentFlags().BoolVar(&config.Config.Secure, "secure", true, "Use secure connection to server")
	RootCmd.PersistentFlags().StringVarP(&config.Config.APIKey, "apiKey", "K", "", "API NetworkKey for connection to server")

	RootCmd.AddCommand(contract.ContractCmd)
	RootCmd.AddCommand(chain.ChainCmd)
	RootCmd.AddCommand(event.EventCmd)
	RootCmd.AddCommand(fungibles.FungibleCmd)

	// Under construction
	RootCmd.AddCommand(stats.StatsCmd)
}
