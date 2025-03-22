package swap

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/heimdahl-xyz/heimdahl-cli/format"
	"github.com/spf13/cobra"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var page int
var perPage int
var formatF string

// Swap represents a token swap transaction
type Swap struct {
	ChainName           string   `json:"chain_name"`
	TxHash              string   `json:"tx_hash"`
	Timestamp           int64    `json:"timestamp"`
	Token1Address       string   `json:"token1_address"`
	Token1Symbol        string   `json:"token1_symbol"`
	Token1Decimals      int      `json:"token1_decimals"`
	Token2Address       string   `json:"token2_address"`
	Token2Symbol        string   `json:"token2_symbol"`
	Token2Decimals      int      `json:"token2_decimals"`
	Token1Sender        string   `json:"token1_sender"`
	Token2Sender        string   `json:"token2_sender"`
	Token1Amount        *big.Int `json:"token1_amount"`
	Token2Amount        *big.Int `json:"token2_amount"`
	PriceToken1InToken2 *big.Int `json:"price_token1_in_token2,string"`
	PriceToken2InToken1 *big.Int `json:"price_token2_in_token1,string"`
}

// SwapData represents the structure of the JSON data
type SwapData struct {
	Meta struct {
		Timestamp int64    `json:"timestamp"`
		Chains    []string `json:"chains"`
		Tokens    []string `json:"tokens"`
		Page      int      `json:"page"`
		PerPage   int      `json:"per_page"`
		Total     int      `json:"total"`
	} `json:"meta"`
	Swaps []Swap `json:"swaps"`
}

// formatTimestamp converts Unix timestamp to human-readable time
func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// RenderSwapsTable renders the token swaps as a table
func RenderSwapsTable(jsonData []byte) error {
	log.Printf("%s", jsonData)
	var swapData SwapData
	err := json.Unmarshal(jsonData, &swapData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Calculate column widths
	timeWidth := 19   // "YYYY-MM-DD HH:MM:SS"
	chainWidth := 10  // "ethereum", "polygon", etc.
	tokenWidth := 6   // "USDT", "WETH", etc.
	amountWidth := 15 // Token amounts
	txWidth := 66     // Full transaction hashes

	// Print table header with metadata
	fmt.Printf("Token Swaps (%d found)\n", swapData.Meta.Total)
	fmt.Printf("Chains: %s\n", strings.Join(swapData.Meta.Chains, ", "))
	fmt.Printf("Tokens: %s\n", strings.Join(swapData.Meta.Tokens, ", "))
	fmt.Printf("Page: %d (showing %d per page)\n\n", swapData.Meta.Page+1, swapData.Meta.PerPage)

	// Define the divider line
	dividerLine := fmt.Sprintf("+-%s-+-%s-+-%s-+-%s-+-%s-+-%s-+",
		strings.Repeat("-", timeWidth),
		strings.Repeat("-", chainWidth),
		strings.Repeat("-", txWidth),
		strings.Repeat("-", tokenWidth),
		strings.Repeat("-", tokenWidth),
		strings.Repeat("-", amountWidth*2+3)) // +3 for the "for" text

	// Print table header
	fmt.Println(dividerLine)
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		timeWidth, "Time",
		chainWidth, "Chain",
		txWidth, "Transaction Hash",
		tokenWidth, "From",
		tokenWidth, "To",
		amountWidth*2+3, "Amount")
	fmt.Println(dividerLine)

	// Print each swap
	for _, swap := range swapData.Swaps {
		// Format amounts

		amount1 := format.FormatAmountBigInt(swap.Token1Amount, uint8(swap.Token1Decimals))
		amount2 := format.FormatAmountBigInt(swap.Token2Amount, uint8(swap.Token2Decimals))

		// Format the combined amount string
		amountStr := fmt.Sprintf("%s %s for %s %s",
			amount1, swap.Token1Symbol,
			amount2, swap.Token2Symbol)

		// Print the row
		fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
			timeWidth, formatTimestamp(swap.Timestamp),
			chainWidth, swap.ChainName,
			txWidth, swap.TxHash,
			tokenWidth, swap.Token1Symbol,
			tokenWidth, swap.Token2Symbol,
			amountWidth*2+3, amountStr)
	}

	// Close the table
	fmt.Println(dividerLine)

	return nil
}

// RenderSwapsCSV writes the swaps data to stdout in CSV format
func RenderSwapsCSV(jsonData []byte) error {
	var swapData SwapData
	err := json.Unmarshal(jsonData, &swapData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Create a CSV writer
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Define headers
	headers := []string{
		"Timestamp",
		"Time",
		"Chain",
		"TX Hash",
		"From Token",
		"From Token Address",
		"From Amount Raw",
		"From Amount Formatted",
		"To Token",
		"To Token Address",
		"To Amount Raw",
		"To Amount Formatted",
		"From Sender",
		"To Sender",
		"Price Token1 In Token2",
		"Price Token2 In Token1",
	}

	// Write header row
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write data rows
	for _, swap := range swapData.Swaps {
		// Format amounts
		amount1Formatted := format.FormatAmountBigInt(swap.Token1Amount, uint8(swap.Token1Decimals))
		amount2Formatted := format.FormatAmountBigInt(swap.Token2Amount, uint8(swap.Token2Decimals))

		// Generate price strings
		var price1In2Str, price2In1Str string
		if swap.PriceToken1InToken2 != nil {
			price1In2Str = swap.PriceToken1InToken2.String()
		}
		if swap.PriceToken2InToken1 != nil {
			price2In1Str = swap.PriceToken2InToken1.String()
		}

		// Create row
		row := []string{
			strconv.FormatInt(swap.Timestamp, 10),
			formatTimestamp(swap.Timestamp),
			swap.ChainName,
			swap.TxHash,
			swap.Token1Symbol,
			swap.Token1Address,
			swap.Token1Amount.String(),
			amount1Formatted,
			swap.Token2Symbol,
			swap.Token2Address,
			swap.Token2Amount.String(),
			amount2Formatted,
			swap.Token1Sender,
			swap.Token2Sender,
			price1In2Str,
			price2In1Str,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	// Add metadata as comments
	metaComments := [][]string{
		{"# Metadata:"},
		{"# Export Time", time.Now().Format("2006-01-02 15:04:05")},
		{"# Data Timestamp", formatTimestamp(swapData.Meta.Timestamp)},
		{"# Total Swaps", strconv.Itoa(swapData.Meta.Total)},
		{"# Page", strconv.Itoa(swapData.Meta.Page + 1)},
		{"# Per Page", strconv.Itoa(swapData.Meta.PerPage)},
		{"# Chains", strings.Join(swapData.Meta.Chains, ", ")},
		{"# Tokens", strings.Join(swapData.Meta.Tokens, ", ")},
	}

	for _, comment := range metaComments {
		if err := writer.Write(comment); err != nil {
			return fmt.Errorf("error writing CSV metadata: %v", err)
		}
	}

	return nil
}

// ListCmd represents the listen command
var ListCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "list  swaps for fungible tokens by pattern",
	Long: `List fungible token swaps 
	Arguments:
	  pattern - search pattern (required) (eg. ethereum.mainnet.usdt.weth.all)`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		// Prepare the WebSocket URL
		hurl := fmt.Sprintf("%s/v1/swaps/list/%s?page=%d&pageSize=%d", config.GetHost(), pattern, page, perPage)

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

		switch formatF {
		case "table":
			err = RenderSwapsTable(b)
		case "csv":
			err = RenderSwapsCSV(b)
		case "json":
			fmt.Println(string(b))
		}

		if err != nil {
			log.Println("failed to render response %s", err)
			return
		}
	},
}

func init() {
	ListCmd.Flags().IntVar(&page, "page", 0, "Page of returned results")
	ListCmd.Flags().IntVar(&perPage, "perPage", 20, "Size of page")
	ListCmd.Flags().StringVar(&formatF, "format", "table", "Output format (table,csv,json)")
}
