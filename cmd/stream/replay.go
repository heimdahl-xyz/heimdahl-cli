package stream

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type EventDetails struct {
	Details []map[string]interface{} `json:"events"`
}

// Assuming we have event as map[string]interface{}
func renderReplayEventTable(event map[string]interface{}) {

	// Format known fields
	blkn := event["blockNumber"].(float64)
	blockNum := strconv.FormatInt(int64(blkn), 10)
	blockHash := fmt.Sprintf("%-15s", event["blockHash"].(string))
	timestamp := event["blockTimestamp"].(string)
	txindex := int64(event["transactionIndex"].(float64))

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
	fmt.Printf("%s | %s | %s | %d | %s\n",
		blockNum, blockHash, timestamp, txindex, strings.Join(eventData, ", "))
}

func isReplayMetaField(field string) bool {
	metaFields := map[string]bool{
		"blockNumber":      true,
		"blockHash":        true,
		"blockTimestamp":   true,
		"timestamp":        true,
		"transactionHash":  true,
		"transactionIndex": true,
	}
	return metaFields[field]
}

// ListenCmd represents the listen command
var ReplayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Replay events for single or multiple contracts",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		event, _ := cmd.Flags().GetString("event")

		if address == "" {
			log.Fatal("Address must be provided")
		}

		// Prepare the WebSocket URL
		httpUrl := fmt.Sprintf("%s/v1/events?address=%s&event=%s", config.GetHost(), address, event)

		headers := make(http.Header)

		headers.Set("X-API-Key", config.GetApiKey())
		headers.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Get(httpUrl)
		if err != nil {
			log.Fatalf("unable to retrieve events %s", err)
		}

		var details EventDetails
		err = json.NewDecoder(resp.Body).Decode(&details)
		if err != nil {
			log.Fatalf("unable to parse details %s", err)
		}

		theaders := []string{"BLOCK#", "BLOCK_HASH", "TIMESTAMP", "CONTRACT", "TRANSACTION_HASH", "EVENT_DATA"}

		//// Print header row
		fmt.Printf("%-10s | %-65s | %-8s | %-15s | %-19s | %-15s\n",
			theaders[0],
			theaders[1],
			theaders[2],
			theaders[3],
			theaders[4],
			theaders[5])

		fmt.Println(strings.Repeat("-", 100))

		for i := range details.Details {
			renderReplayEventTable(details.Details[i])
		}
	},
}

func init() {
	// Define flags
	ReplayCmd.Flags().StringP("address", "a", "", "WebSocket server address")
	ReplayCmd.Flags().StringP("event", "e", "", "Event to replay")
}
