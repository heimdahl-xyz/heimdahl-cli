package event

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Assuming we have event as map[string]interface{}
func renderEventTable(event map[string]interface{}) {
	//log.Printf("%+v", event)

	//// Format known fields
	blkn := event["blockNumber"].(float64)
	blockNum := strconv.FormatInt(int64(blkn), 10)
	blockHash := fmt.Sprintf("%-15s", event["blockHash"].(string))
	timestamp := time.Unix(int64(event["blockTimestamp"].(float64)), 0).Format("2006-01-02 15:04:05")
	//contract := fmt.Sprintf("%-15s", event["contractAddress"].(string))

	var eventData []string
	for k, v := range event {
		// Skip already used fields
		if isMetaField(k) {
			continue
		}
		eventData = append(eventData, fmt.Sprintf("%s: %v", k, v))
	}

	fmt.Printf("| %s | %s | %s | %s\n",
		blockNum, blockHash, timestamp, strings.Join(eventData, ", "))
}

func isMetaField(field string) bool {
	metaFields := map[string]bool{
		"chain":            true,
		"network":          true,
		"blockNumber":      true,
		"blockHash":        true,
		"blockTimestamp":   true,
		"contractAddress":  true,
		"timestamp":        true,
		"transactionHash":  true,
		"transactionIndex": true,
	}
	return metaFields[field]
}

// SubscribeCmd represents the listen command
var SubscribeCmd = &cobra.Command{
	Use:   "subscribe [address] [event]",
	Short: "subscribe to realtime events for contract",
	Long: `Subscribe to realtime events for contract 
	Arguments:
	  address - The contract address (required) (eg. 0xdAC17F958D2ee523a2206206994597C13D831ec7)
	  event-name   - Name of the event (eg. Transfer),`,
	Args: cobra.ExactArgs(2), // Expect exactly 2 arguments

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			return
		}

		address := args[0]
		event := args[1]

		// Prepare the WebSocket URL
		wsURL := fmt.Sprintf("%s/v1/%s/listen/%s/%s?api_key=%s", config.GetWsHost(), chain, address, event, config.GetApiKey())

		log.Println("Connecting to WebSocket: ", wsURL)

		//log.Println(wsURL)
		headers := make(http.Header)

		headers.Set("Content-Type", "application/json")

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Fatal("Error connecting to WebSocket:", err)
		}
		defer conn.Close()

		// Define headers
		theaders := []string{"BLOCK#", "BLOCK_HASH", "TIMESTAMP", "CONTRACT", "TRANSACTION_HASH", "EVENT_DATA"}

		// Print header row
		fmt.Printf("%-8s | %-15s | %-19s | %-15s | %-15s | %s\n",
			theaders[0],
			theaders[1],
			theaders[2],
			theaders[3],
			theaders[4],
			theaders[5])

		// Print separator
		fmt.Println(strings.Repeat("-", 120))

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
			//log.Printf("%+v", event)
			renderEventTable(event)
		}
	},
}

func init() {
	SubscribeCmd.Flags().StringVarP(&chain, "chain", "c", "ethereum", "Blockchain type  (eg. ethereum, required)")
	SubscribeCmd.Flags().StringVarP(&network, "network", "w", "mainnet", "Blockchain network (eg. mainnet, required)")
}
