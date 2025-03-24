package subscription

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List subscriptions",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/chain", config.GetHost()) // Use the global host variable

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+config.GetApiKey())
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to list event listeners: %s\n", resp.Status)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var chainInfos []ChainInfo
		err = json.Unmarshal(body, &chainInfos)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		fmt.Printf("\n%-10s %-10s %-8s\n", "CHAIN", "NETWORK", "CHAIN ID")
		fmt.Println("--------------------------------")
		for _, contractInfo := range chainInfos {
			fmt.Printf("%-10s %-10s %-8d\n",
				contractInfo.Chain,
				contractInfo.Network,
				contractInfo.ChainID)
		}
		fmt.Println()
	},
}
