package cmd

import (
	"fmt"
	chain "github.com/heimdahl-xyz/heimdahl-cli/cmd/chains"
	"github.com/heimdahl-xyz/heimdahl-cli/cmd/contract"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "heimdahl",
	Short: "A CLI client for interacting with the Heimdahl event listener API",
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
	RootCmd.PersistentFlags().StringVarP(&config.Config.APIKey, "apiKey", "K", "", "API Key for connection to server")

	RootCmd.AddCommand(contract.ContractCmd)
	RootCmd.AddCommand(chain.ChainCmd)
	RootCmd.AddCommand(ListenCmd)
	RootCmd.AddCommand(ReplayCmd)
}
