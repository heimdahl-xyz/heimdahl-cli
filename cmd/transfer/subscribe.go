package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/heimdahl-xyz/heimdahl-cli/lib"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func renderTokenTransferAsTableRow(transfer lib.FungibleTokenTransfer) {
	fmt.Printf("| %-15d | %-20s | %-20s | %-20s | %-20s | %-15s | %-15s | %-10s | %-10s | %-25s | %-10d | %-10d |\n",
		transfer.Timestamp,
		transfer.FromAddress,
		transfer.FromOwner,
		transfer.ToAddress,
		transfer.ToOwner,
		transfer.Amount.String(),
		transfer.TokenAddress,
		transfer.Symbol,
		transfer.Chain,
		transfer.Network,
		transfer.TxHash,
		transfer.Decimals,
		transfer.Position)
	log.Println("---------------------------------------------------------------------------")

}

// SubscribeCmd represents the listen command
var SubscribeCmd = &cobra.Command{
	Use:   "subscribe [pattern]",
	Short: "subscribe to realtime transfers for fungibe tokens by pattern",
	Long: `Subscribe to realtime events for contract 
	Arguments:
	  pattern - search pattern (required) (eg. ethereum.mainnet.0x1234.0x5677.whale)`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		// Prepare the WebSocket URL
		wsURL := fmt.Sprintf("%s/v1/transfer/stream/%s?api_key=%s", config.GetWsHost(), pattern, config.GetApiKey())

		//log.Println(wsURL)
		headers := make(http.Header)

		headers.Set("Content-Type", "application/json")

		signalChannel := make(chan os.Signal, 1)
		// Notify when SIGINT (Ctrl+C) or SIGTERM signal is received
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

		conn, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
		if err != nil {
			log.Fatal("Error connecting to WebSocket:", err)
		}
		defer conn.Close()

		// separate goroutine to listen for signals
		go func() {
			<-signalChannel
			os.Exit(0)
		}()

		// Format and print the struct fields as a table row

		fmt.Printf("| %-15s | %-20s | %-20s | %-20s | %-20s | %-15s | %-15s | %-10s | %-10s | %-25s | %-10d | %-10d |\n",
			"Timestamp",
			"From Address",
			"From Owner",
			"To Address",
			"To Owner",
			"Amount",
			"Token Address",
			"Symbol",
			"Chain",
			"Network",
			"Tx Hash",
			"Decimals",
			"Position")

		// Listen for messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}
			//log.Println(string(message))
			var transfer lib.FungibleTokenTransfer
			err = json.Unmarshal(message, &transfer)
			if err != nil {
				log.Printf("raw message %s", message)
				log.Println("Error unmarshalling message:", err)
				return
			}
			PrintTransfer(&transfer)
		}

	},
}
