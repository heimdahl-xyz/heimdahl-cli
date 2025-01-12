package event

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

// ListenCmd represents the listen command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List events for single contract",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			cmd.Help()
			return
		}

		chain := args[0]
		address := args[1]
		event := args[2]

		page, _ := cmd.Flags().GetInt32("page")
		perpage, _ := cmd.Flags().GetInt32("perPage")

		// Prepare the WebSocket URL
		httpUrl := fmt.Sprintf("%s/v1/%s/events/%s/%s?page=%d&per_page=%d", config.GetHost(), chain, address, event, page, perpage)

		headers := make(http.Header)

		headers.Set("Authorization", "Bearer "+config.GetApiKey())
		headers.Set("Content-Type", "application/json")

		//log.Println("Requesting events from", httpUrl)

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

		theaders := []string{
			"BLOCK#",
			"TIMESTAMP",
			"TRANSACTION_HASH",
			"EVENT_DATA"}

		//// Print header row
		fmt.Printf("%-10s | %-15s | %-65s | %-15s\n",
			theaders[0],
			theaders[1],
			theaders[2],
			theaders[3],
		)

		fmt.Println(strings.Repeat("-", 100))

		for i := range details.Details {
			renderReplayEventTable(details.Details[i])
		}
	},
}

func init() {
	ListCmd.Flags().Int32P("page", "p", 0, "Page to replay")
	ListCmd.Flags().Int32P("perPage", "n", 20, "Events per page")
}
