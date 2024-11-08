package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen to a stream of events from a WebSocket server",
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		events, _ := cmd.Flags().GetString("events")

		if address == "" {
			log.Fatal("Address must be provided")
		}

		// Prepare the WebSocket URL
		wsURL := fmt.Sprintf("%s/v1/ws-listen?address=%s&events=%s", getWsHost(), address, url.QueryEscape(events))

		headers := make(http.Header)

		headers.Set("X-API-Key", getApiKey())
		headers.Set("Content-Type", "application/json")

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Fatal("Error connecting to WebSocket:", err)
		}
		defer conn.Close()

		log.Println("Connected to WebSocket:", wsURL)

		// Listen for messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			log.Printf("%s", message)
		}
	},
}

func init() {

	// Define flags
	listenCmd.Flags().StringP("address", "a", "", "WebSocket server address")
	listenCmd.Flags().StringP("events", "e", "", "Comma-separated list of events to subscribe to")
}
