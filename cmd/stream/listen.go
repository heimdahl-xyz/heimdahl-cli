package stream

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// Assuming we have event as map[string]interface{}
func renderEventTable(event map[string]interface{}) {

	// Format known fields
	chain := fmt.Sprintf("%-10s", event["chain"])
	network := fmt.Sprintf("%-10s", event["network"])
	blkn := event["blockNumber"].(float64)
	blockNum := strconv.FormatInt(int64(blkn), 10)
	blockHash := fmt.Sprintf("%-15s", event["blockHash"].(string))
	timestamp := time.Unix(int64(event["blockTimestamp"].(float64)), 0).Format("2006-01-02 15:04:05")
	contract := fmt.Sprintf("%-15s", event["contractAddress"].(string))

	// Collect remaining fields for event data
	var eventData []string
	for k, v := range event {
		// Skip already used fields
		if isMetaField(k) {
			continue
		}
		eventData = append(eventData, fmt.Sprintf("%s: %v", k, v))
	}

	// Print data row
	fmt.Printf("%s | %s | %s | %s | %s | %s | %s\n",
		chain, network, blockNum, blockHash, timestamp, contract, strings.Join(eventData, ", "))
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

func isMetaField(field string) bool {
	metaFields := map[string]bool{
		"chain":           true,
		"network":         true,
		"blockNumber":     true,
		"blockHash":       true,
		"blockTimestamp":  true,
		"contractAddress": true,
		"timestamp":       true,
		"transactionHash": true,
	}
	return metaFields[field]
}

// ListenCmd represents the listen command
var ListenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen to events contract",
	Run: func(cmd *cobra.Command, args []string) {
		addresses, _ := cmd.Flags().GetString("addresses")
		event, _ := cmd.Flags().GetString("event")

		if addresses == "" {
			log.Fatal("Address must be provided")
		}
		// Prepare the WebSocket URL
		wsURL := fmt.Sprintf("%s/v1/listen?addresses=%s&event=%s", config.GetWsHost(), addresses, url.QueryEscape(event))

		//log.Println(wsURL)
		headers := make(http.Header)

		headers.Set("X-API-Key", config.GetApiKey())
		headers.Set("Content-Type", "application/json")

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Fatal("Error connecting to WebSocket:", err)
		}
		defer conn.Close()

		// Define headers
		theaders := []string{"CHAIN", "NETWORK", "BLOCK#", "BLOCK_HASH", "TIMESTAMP", "CONTRACT", "TRANSACTION_HASH", "EVENT_DATA"}

		// Print header row
		fmt.Printf("%-10s | %-10s | %-8s | %-15s | %-19s | %-15s | %-15s | %s\n",
			theaders[0],
			theaders[1],
			theaders[2],
			theaders[3],
			theaders[4],
			theaders[5],
			theaders[6],
			theaders[7])

		// Print separator
		fmt.Println(strings.Repeat("-", 100))

		// Listen for messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			var event map[string]interface{}

			err = json.Unmarshal(message, &event)
			if err != nil {
				log.Println("Error unmarshalling message:", err)
				return
			}

			renderEventTable(event)
		}
	},
}

func init() {

	// Define flags
	ListenCmd.Flags().StringP("addresses", "a", "", "Contract address to listen to")
	ListenCmd.Flags().StringP("event", "e", "", "Comma-separated list of events to subscribe to")
}
