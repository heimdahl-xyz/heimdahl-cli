package chain

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type ChainInfo struct {
	Chain   string `json:"chain_name"`
	Network string `json:"chain_network"`
	ChainID int    `json:"chain_id"`
}

var (
	chain   string
	network string
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show chain information",
	Run: func(cmd *cobra.Command, args []string) {
		// change the endpoint
		url := fmt.Sprintf("%s/v1/chains/%s/%s", config.GetHost(), chain, network) // Use the global host variable
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Println("Error making GET request:", err)
			return
		}

		req.Header.Set("X-API-Key", config.GetApiKey())
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

		var contractInfo ChainInfo
		err = json.Unmarshal(body, &contractInfo)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		fmt.Printf("Chain: %s\nNetwork: %s\nChain ID:%d\n ", contractInfo.Chain, contractInfo.Network, contractInfo.ChainID)
	},
}

func init() {
	ShowCmd.Flags().StringVarP(&chain, "chain", "c", "", "Chain name (required)")
	ShowCmd.Flags().StringVarP(&network, "network", "n", "mainnet", "Chain network  (eg. mainnet, sepolia)")
	_ = ShowCmd.MarkFlagRequired("chain")
}
