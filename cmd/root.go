package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	host   string // Global host variable
	secure bool
	apiKey string
)

var rootCmd = &cobra.Command{
	Use:   "heimdahl",
	Short: "A CLI client for interacting with the Heimdahl event listener API",
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream subcommands",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getHost() string {
	if secure {
		return "https://" + host
	} else {
		return "http://" + host
	}
}

func getWsHost() string {
	if secure {
		return "wss://" + host
	} else {
		return "ws://" + host
	}
}

func getApiKey()string {
	apk := os.Getenv("HEIMDAHL_API_KEY")
	if apk == "" {
		apk = apiKey
	}
	return apk
}


func init() {
	// Add a global flag for host
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "api.heimdahl.xyz", "Host URL for the API server")
	rootCmd.PersistentFlags().BoolVar(&secure, "secure", true,  "Use secure connection to server")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apiKey", "K", "", "API Key for connection to server")

	// Add subcommands here
	streamCmd.AddCommand(addCmd)
	streamCmd.AddCommand(getCmd)
	streamCmd.AddCommand(listCmd)
	streamCmd.AddCommand(listenCmd)

	rootCmd.AddCommand(streamCmd)
}
