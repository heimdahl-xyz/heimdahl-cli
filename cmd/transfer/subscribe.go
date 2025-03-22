package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/heimdahl-xyz/heimdahl-cli/lib"
	"github.com/spf13/cobra"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// formatAmount formats token amount according to its decimals
func formatAmountBigInt(amount *big.Int, decimals uint8) string {
	if amount == nil {
		return "0"
	}

	// Clone the amount to avoid modifying the original
	result := new(big.Int).Set(amount)

	// If decimals is 0, just return the amount as string
	if decimals == 0 {
		return result.String()
	}

	// Calculate the divisor (10^decimals)
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// Calculate the integer part
	intPart := new(big.Int).Div(result, divisor)

	// Calculate the fractional part
	fracPart := new(big.Int).Mod(result, divisor)

	// Convert to string with proper padding
	fracStr := fracPart.String()

	// Pad with leading zeros if necessary
	for uint8(len(fracStr)) < decimals {
		fracStr = "0" + fracStr
	}

	// Trim trailing zeros
	fracStr = strings.TrimRight(fracStr, "0")

	// Format the final string
	if fracStr == "" {
		return intPart.String()
	}
	return intPart.String() + "." + fracStr
}

// PrintTransfer prints a single transfer in a human-readable format
func PrintTransfer(transfer *lib.FungibleTokenTransfer) {
	// Create a horizontal line with a longer length to accommodate full addresses
	horizLine := strings.Repeat("─", 120)

	fmt.Println("┌" + horizLine + "┐")

	// Format: Chain.Network → Symbol Transfer
	chainStr := transfer.Chain
	if transfer.Network != "" && transfer.Network != "mainnet" {
		chainStr += "." + transfer.Network
	}
	fmt.Printf("│ \033[1m%s → %s Transfer\033[0m\n", chainStr, transfer.Symbol)

	// Time and Transaction
	fmt.Printf("│ \033[90mTimestamp:\033[0m %s\n", formatTimestamp(transfer.Timestamp))
	fmt.Printf("│ \033[90mTX Hash:  \033[0m %s\n", transfer.TxHash)

	// Add a separator
	fmt.Println("│ " + strings.Repeat("─", 118))

	// From → To with owners if available
	fmt.Printf("│ \033[90mFrom:     \033[0m %s", transfer.FromAddress)
	if transfer.FromOwner != "" && transfer.FromOwner != transfer.FromAddress {
		fmt.Printf(" \033[90m(Owner: %s)\033[0m", transfer.FromOwner)
	}
	fmt.Println()

	fmt.Printf("│ \033[90mTo:       \033[0m %s", transfer.ToAddress)
	if transfer.ToOwner != "" && transfer.ToOwner != transfer.ToAddress {
		fmt.Printf(" \033[90m(Owner: %s)\033[0m", transfer.ToOwner)
	}
	fmt.Println()

	// Amount with symbol
	formattedAmount := formatAmountBigInt(transfer.Amount, transfer.Decimals)
	fmt.Printf("│ \033[90mAmount:   \033[0m \033[1m%s %s\033[0m\n", formattedAmount, transfer.Symbol)

	// Token address
	fmt.Printf("│ \033[90mToken:    \033[0m %s\n", transfer.TokenAddress)

	// Position
	fmt.Printf("│ \033[90mPosition: \033[0m %d\n", transfer.Position)

	fmt.Println("└" + horizLine + "┘")
}

// PrintTransfers prints multiple transfers with a header
func PrintTransfers(transfers []*lib.FungibleTokenTransfer) {
	if len(transfers) == 0 {
		fmt.Println("No transfers found.")
		return
	}

	fmt.Printf("Found %d transfers:\n\n", len(transfers))

	for i, transfer := range transfers {
		PrintTransfer(transfer)

		// Add a newline between transfers, but not after the last one
		if i < len(transfers)-1 {
			fmt.Println()
		}
	}
}

// SubscribeCmd represents the listen command
var SubscribeCmd = &cobra.Command{
	Use:   "subscribe [pattern]",
	Short: "subscribe to realtime transfer for fungibe tokens by pattern",
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
