package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type ContractInfo struct {
	Network         string `json:"network"`
	ContractName    string `json:"contract_name"`
	ContractAddress string `json:"contract_address"`
	ABI             string `json:"-"`
}

var address string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an event listener by address",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/event-listeners/%s", host, address) // Use the global host variable
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}

		req.Header.Set("X-API-Key", apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error performing request:", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to get event listener: %s\n", resp.Status)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var contractInfo ContractInfo
		err = json.Unmarshal(body, &contractInfo)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		fmt.Printf("Network: %s\nContract Name: %s\nContract Address: %s\n", contractInfo.Network, contractInfo.ContractName, contractInfo.ContractAddress)
	},
}

func init() {
	getCmd.Flags().StringVarP(&address, "address", "a", "", "Contract address (required)")
	_ = getCmd.MarkFlagRequired("address")
}
