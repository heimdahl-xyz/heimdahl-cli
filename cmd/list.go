package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all event listeners",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/event-listeners", host) // Use the global host variable
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to list event listeners: %s\n", resp.Status)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
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

func init() {
	// No flags needed for the list command
}
