package event

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

var (
	chain   string
	network string
)

type EventMeta struct {
	Chain   string `json:"chain"`
	ChainID int    `json:"chain_id"`
	Address string `json:"addresss"` // Note: "addresss" has a typo; update field name accordingly if needed.
	Event   string `json:"event"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Total   int    `json:"total"`
}

type EventData struct {
	Chain   string `json:"chain"`
	ChainID int    `json:"chain_id"`
	Address string `json:"addresss"` // Note: "addresss" has a typo; update field name accordingly if needed.
	Event   string `json:"event"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Total   int    `json:"total"`
}

type EventDetails struct {
	Meta    EventMeta                `json:"meta"`
	Details []map[string]interface{} `json:"events"`
}

// Assuming we have event as map[string]interface{}
func renderReplayEventTable(event map[string]interface{}) {
	// Format known fields
	blkn := event["blockNumber"].(float64)
	blockNum := strconv.FormatInt(int64(blkn), 10)
	transactionHash := fmt.Sprintf("%-15s", event["transactionHash"].(string))
	timestamp := event["blockTimestamp"].(string)

	// Collect remaining fields for event data
	var eventData []string
	for k, v := range event {
		// Skip already used fields
		if isReplayMetaField(k) {
			continue
		}
		eventData = append(eventData, fmt.Sprintf("%s: %v", k, v))
	}

	// Print data row
	fmt.Printf("%s | %s | %s | %s\n",
		blockNum, timestamp, transactionHash, strings.Join(eventData, ", "))
}

func isReplayMetaField(field string) bool {
	metaFields := map[string]bool{
		"blockNumber":     true,
		"blockTimestamp":  true,
		"timestamp":       true,
		"transactionHash": true,
	}
	return metaFields[field]
}

// PrintEventDetails prints EventDetails in a clean, console-friendly format.
// Keys in each event map are sorted alphabetically.
func PrintEventDetails(e EventDetails) {
	separator := strings.Repeat("-", 75)

	// --- META SECTION ---
	fmt.Println(separator)
	fmt.Println("META")
	fmt.Println(separator)

	if len(e.Details) < 1 {
		return
	}

	// NOTE: adapt these fields to your actual EventMeta struct.
	fmt.Printf("Timestamp : %s\n", e.Details[0]["blockTimestamp"])
	fmt.Printf("Page      : %d\n", e.Meta.Page)
	fmt.Printf("Per Page  : %d\n", e.Meta.PerPage)
	fmt.Printf("Total     : %d\n", e.Meta.Total)

	// --- EVENTS SECTION ---
	for _, event := range e.Details {
		fmt.Println(separator)

		// Collect keys
		keys := make([]string, 0, len(event))
		for k := range event {
			keys = append(keys, k)
		}

		// Sort alphabetically
		sort.Strings(keys)

		// Print in order
		for _, k := range keys {
			fmt.Printf("%-10s: %v\n", k, event[k])
		}
	}

	fmt.Println(separator)
}

// SubscribeCmd represents the listen command
var ListCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "List events for contract",
	Long: `List collected events for contract 
Arguments:
	pattern - The search pattern (required) (eg. ethereum.mainnet.0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2.Transfer)`,
	Args: cobra.ExactArgs(1), // Expect exactly 2 arguments

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		page, _ := cmd.Flags().GetInt("page")
		perpage, _ := cmd.Flags().GetInt("perPage")

		// Prepare the WebSocket URL
		httpUrl := fmt.Sprintf("%s/v1/events/list/%s?page=%d&pageSize=%d", config.GetHost(), pattern, page, perpage)
		//log.Println("url ", httpUrl)
		headers := make(http.Header)

		headers.Set("Authorization", "Bearer "+config.GetApiKey())
		headers.Set("Content-Type", "application/json")
		req, _ := http.NewRequest(http.MethodGet, httpUrl, nil)

		req.Header = headers
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("unable to retrieve events %s", err)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("unable to read body %s", err)
		}

		var details EventDetails
		err = json.Unmarshal(b, &details)
		if err != nil {
			log.Fatalf("unable to parse details %s", err)
		}

		PrintEventDetails(details)
	},
}

func init() {
	ListCmd.Flags().StringVarP(&chain, "chain", "c", "ethereum", "Blockchain type  (eg. ethereum, required)")
	ListCmd.Flags().StringVarP(&network, "network", "w", "mainnet", "Blockchain network (eg. mainnet, required)")
	ListCmd.Flags().IntP("page", "p", 0, "Page to replay")
	ListCmd.Flags().IntP("perPage", "l", 20, "Events per page")
}
