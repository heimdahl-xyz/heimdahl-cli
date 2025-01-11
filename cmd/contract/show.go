package contract

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type ContractInfo struct {
	Chain           string `json:"chain"`
	Network         string `json:"network"`
	ContractName    string `json:"contract_name"`
	ContractAddress string `json:"contract_address"`
	ABI             string `json:"-"`
}

var address string

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show contract by address",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("%s/v1/contracts/%s", config.GetHost(), address) // Use the global host variable
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+config.GetApiKey())
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
	ShowCmd.Flags().StringVarP(&address, "address", "a", "", "Contract address (required)")
	_ = ShowCmd.MarkFlagRequired("address")
}
