package contract

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all contracts",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/contracts?chain=%s&network=%s", config.GetHost(), chain, network) // Use the global host variable

		req, err := http.NewRequest("GET", url, nil)
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

		var contractInfos []ContractInfo
		err = json.Unmarshal(body, &contractInfos)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		if len(contractInfos) == 0 {
			fmt.Printf("Could not find indexed contracts for %s network %s\n", chain, network)
			return
		}

		// Then data rows
		for _, contractInfo := range contractInfos {
			fmt.Printf("Chain:            %s\n", contractInfo.Chain)
			fmt.Printf("Network:          %s\n", contractInfo.Network)
			fmt.Printf("Contract Identifier:    %s\n", contractInfo.ContractName)
			fmt.Printf("Contract Address: %s\n", contractInfo.ContractAddress)
			fmt.Printf("Events:\n")
			for _, event := range strings.Split(contractInfo.Events, ",") {
				fmt.Printf("  - %s\n", strings.TrimSpace(event))
			}
			fmt.Println(strings.Repeat("-", 80)) // Add a separator for better readability
		}
		fmt.Println()
	},
}

func init() {
	ListCmd.Flags().StringVarP(&chain, "chain", "c", "ethereum", "Blockchain name (eg. ethereum, required)")
	ListCmd.Flags().StringVarP(&network, "network", "n", "mainnet", "Blockchain network (eg. mainnet, required)")
}
