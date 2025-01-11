package stream

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"strings"
)

// CrossListenCmd represents the cross-listen command
var CrossListenCmd = &cobra.Command{
	Use:   "cross-listen",
	Short: "Cross listen to events on multiple chains",
	Run: func(cmd *cobra.Command, args []string) {
		contract, _ := cmd.Flags().GetString("contract")
		event, _ := cmd.Flags().GetString("event")
		chains, _ := cmd.Flags().GetString("chains")

		if contract == "" {
			log.Fatal("contract must be provided")
		}

		if chains == "" || len(strings.Split(chains, ",")) == 0 {
			log.Fatal("chains must be provided")
		}

		if event == "" {
			log.Fatal("event must be provided")
		}

		// Prepare the WebSocket URL
		wsURL := fmt.Sprintf("%s/v1/cross-listen?contract=%s&event=%s&chains=%s",
			config.GetWsHost(),
			contract,
			event,
			chains)

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

			renderEventTable(event)
		}
	},
}

func init() {

	// Define flags
	CrossListenCmd.Flags().StringP("contract", "a", "", "Contract name to listen to")
	CrossListenCmd.Flags().StringP("event", "e", "", "Comma-separated list of events to subscribe to")
	CrossListenCmd.Flags().StringP("chains", "c", "", "Comma-separated list of chains to subscribe to")
}
