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
	Events          string `json:"events"`
	ABI             string `json:"-"`
}

var address string

var ShowCmd = &cobra.Command{
	Use:   "show [address]",
	Short: "Show contract by address",
	Long: `Show contract metadata by address. 
		Usage: heimdahl contract show 0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: heimdahl contract show [address]")
			return
		}
		address := args[0]

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

		fmt.Printf("Network: %s\nContract Name: %s\nContract Address: %s\nEvents: %s",
			contractInfo.Network,
			contractInfo.ContractName,
			contractInfo.ContractAddress,
			contractInfo.Events)
	},
}

func init() {

}
