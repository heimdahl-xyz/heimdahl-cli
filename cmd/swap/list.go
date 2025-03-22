package swap

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
	Short: "list  swaps for fungible tokens by pattern",
	Long: `List fungible token swaps 
	Arguments:
	  pattern - search pattern (required) (eg. ethereum.mainnet.usdt.weth.whale)`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		// Prepare the WebSocket URL
		hurl := fmt.Sprintf("%s/v1/swap/list/%s?page=%d&pageSize=%d", config.GetHost(), pattern, page, perPage)

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
	},
}

func init() {
	ListCmd.Flags().IntVar(&page, "page", 0, "Page of returned results")
	ListCmd.Flags().IntVar(&perPage, "perPage", 20, "Size of page")
}
