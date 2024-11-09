package contract

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
	Short: "List all event listener streams",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/contracts", config.GetHost()) // Use the global host variable

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			return
		}

		req.Header.Set("X-API-Key", config.GetApiKey())
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
		// First print header
		fmt.Printf("\n%-10s %-15s %-42s\n",
			"NETWORK",
			"CONTRACT NAME",
			"CONTRACT ADDRESS")
		fmt.Println("-------------------------------------------------------------------------------")

		// Then data rows
		for _, contractInfo := range contractInfos {
			fmt.Printf("%-10s | %-10s | %-15s | %-42s\n",
				contractInfo.Chain,
				contractInfo.Network,
				contractInfo.ContractName,
				contractInfo.ContractAddress)
		}
		fmt.Println()
	},
}
