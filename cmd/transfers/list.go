package transfers

import (
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
)

var page int
var perPage int

// ListCmd represents the listen command
var ListCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "list transfers for fungible tokens by pattern",
	Long: `List fungible token transfers 
	Arguments:
	  pattern - search pattern (required) (eg. ethereum.mainnet.usdt.0x1234.0x5677.whale)`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		// Prepare the WebSocket URL
		hurl := fmt.Sprintf("%s/v1/transfers/list/%s?page=%d&pageSize=%d", config.GetHost(), pattern, page, perPage)

		req, _ := http.NewRequest(http.MethodGet, hurl, nil)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.GetApiKey())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to perform request %s", err)
			return
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("failed to read response %s", err)
			return
		}

		fmt.Printf("%s", b)
		// Format and print the struct fields as a table row
		//fmt.Printf("| %-15s | %-20s | %-20s | %-20s | %-20s | %-15s | %-15s | %-10s | %-10s | %-25s | %-10d | %-10d |\n",
		//	"Timestamp", "From Address", "From Owner", "To Address", "To Owner", "Amount", "Token Address", "Symbol", "Chain", "Network", "Tx Hash", "Decimals", "Position")

		//renderTokenTransferAsTableRow(transfer)

	},
}

func init() {
	ListCmd.Flags().IntVar(&page, "page", 0, "Page of returned results")
	ListCmd.Flags().IntVar(&perPage, "perPage", 20, "Size of page")
}
