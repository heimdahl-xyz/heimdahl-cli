package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	host   string // Global host variable
	wsHost string
	apiKey string
)
var rootCmd = &cobra.Command{
	Use:   "heim-cli",
	Short: "A CLI client for interacting with the Heimdahl event listener API",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add a global flag for host
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "https://api.heimdahl.xyz", "Host URL for the API server")
	rootCmd.PersistentFlags().StringVarP(&wsHost, "wsHost", "W", "wss://api.heimdahl.xyz", "WSHost URL for the API server")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apiKey", "K", "", "API Key for connection to server")

	// Add subcommands here
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(hardhatInitCmd)
	rootCmd.AddCommand(hardhatDeployCmd)

}
