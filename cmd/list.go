package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all event listener streams",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/event-listeners", getHost()) // Use the global host variable

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			return
		}

		req.Header.Set("X-API-Key", getApiKey())
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

		for _, contractInfo := range contractInfos {
			fmt.Printf("Network: %s\nContract Name: %s\nContract Address: %s\n\n", contractInfo.Network, contractInfo.ContractName, contractInfo.ContractAddress)
		}
	},
}
