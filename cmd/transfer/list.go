package transfer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	format2 "github.com/heimdahl-xyz/heimdahl-cli/format"
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
var format string

// Transfer represents a token transfer transaction
type Transfer struct {
	Timestamp    int64    `json:"timestamp"`
	FromAddress  string   `json:"from_address"`
	ToAddress    string   `json:"to_address"`
	Amount       *big.Int `json:"amount"`
	TokenAddress string   `json:"token_address"`
	Symbol       string   `json:"symbol"`
	Chain        string   `json:"chain"`
	Network      string   `json:"network"`
	TxHash       string   `json:"tx_hash"`
	Decimals     int      `json:"decimals"`
	Position     int64    `json:"position"`
}

// TokenResponse represents the structure of the JSON data
type TokenResponse struct {
	Meta struct {
		Timestamp int64    `json:"timestamp"`
		Chains    []string `json:"chains"`
		Tokens    []string `json:"tokens"`
		Page      int      `json:"page"`
		PerPage   int      `json:"per_page"`
		Total     int      `json:"total"`
	} `json:"meta"`
	Transfers []Transfer `json:"transfers"`
}

// formatTimestamp converts Unix timestamp to human-readable time
func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// formatAmount converts *big.Int + decimals into a readable decimal string
func formatAmount(amount *big.Int, decimals int) string {
	if amount == nil {
		return "0"
	}

	str := amount.String()
	l := len(str)

	if decimals == 0 {
		return str
	}

	if l <= decimals {
		return "0." + strings.Repeat("0", decimals-l) + str
	}

	return str[:l-decimals] + "." + str[l-decimals:]
}

func PrintTokenResponse(jsonData []byte) error {
	var resp TokenResponse
	err := json.Unmarshal(jsonData, &resp)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	separator := strings.Repeat("-", 75)

	// --- META SECTION ---
	fmt.Println(separator)
	fmt.Println("META")
	fmt.Println(separator)

	ts := time.Unix(resp.Meta.Timestamp, 0).Format("2006-01-02 15:04:05")

	fmt.Printf("Timestamp : %s\n", ts)
	fmt.Printf("Chains    : %s\n", strings.Join(resp.Meta.Chains, ", "))
	fmt.Printf("Tokens    : %s\n", strings.Join(resp.Meta.Tokens, ", "))
	fmt.Printf("Page      : %d\n", resp.Meta.Page)
	fmt.Printf("Per Page  : %d\n", resp.Meta.PerPage)
	fmt.Printf("Total     : %d\n", resp.Meta.Total)

	// --- TRANSFERS SECTION ---
	for _, t := range resp.Transfers {
		fmt.Println(separator)

		// time conversion
		transferTime := time.Unix(t.Timestamp, 0).Format("2006-01-02 15:04:05")

		// convert amount
		amountStr := formatAmount(t.Amount, t.Decimals)

		fmt.Printf("Time     : %s\n", transferTime)
		fmt.Printf("From     : %s\n", t.FromAddress)
		fmt.Printf("To       : %s\n", t.ToAddress)
		fmt.Printf("Amount   : %s\n", amountStr)
		fmt.Printf("Symbol   : %s\n", t.Symbol)
		fmt.Printf("Chain    : %s\n", t.Chain)
		fmt.Printf("TX Hash  : %s\n", t.TxHash)
	}

	fmt.Println(separator)
	return nil
}

// RenderTransfersToCSV exports token transfer as CSV to stdout
func RenderTransfersToCSV(jsonData []byte) error {
	var tokenData TokenResponse
	err := json.Unmarshal(jsonData, &tokenData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Create a CSV writer that writes to stdout
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header row
	header := []string{
		"Timestamp",
		"Time",
		"From Address",
		"To Address",
		"Amount (Raw)",
		"Amount (Formatted)",
		"Token Symbol",
		"Token Address",
		"Chain",
		"Network",
		"TX Hash",
		"Position",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Write data rows
	for _, transfer := range tokenData.Transfers {
		row := []string{
			strconv.FormatInt(transfer.Timestamp, 10),
			formatTimestamp(transfer.Timestamp),
			transfer.FromAddress,
			transfer.ToAddress,
			format2.FormatAmountBigInt(transfer.Amount, uint8(transfer.Decimals)),
			transfer.Symbol,
			transfer.TokenAddress,
			transfer.Chain,
			transfer.Network,
			transfer.TxHash,
			strconv.FormatInt(transfer.Position, 10),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing CSV row: %v", err)
		}
	}

	// Add metadata at the end as comments
	metaComments := [][]string{
		{"# Metadata:"},
		{"# Export Time", time.Now().Format("2006-01-02 15:04:05")},
		{"# Data Timestamp", formatTimestamp(tokenData.Meta.Timestamp)},
		{"# Total Transfers", strconv.Itoa(tokenData.Meta.Total)},
		{"# Page", strconv.Itoa(tokenData.Meta.Page + 1)},
		{"# Per Page", strconv.Itoa(tokenData.Meta.PerPage)},
		{"# Chains", strings.Join(tokenData.Meta.Chains, ", ")},
		{"# Tokens", strings.Join(tokenData.Meta.Tokens, ", ")},
	}

	for _, comment := range metaComments {
		if err := writer.Write(comment); err != nil {
			return fmt.Errorf("error writing CSV metadata: %v", err)
		}
	}

	return nil
}

// printLine prints a horizontal line for the table
func printLine(width int) {
	fmt.Println(strings.Repeat("-", width))
}

// ListCmd represents the listen command
var ListCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "list transfers for fungible tokens by pattern",
	Long: `List fungible token transfers
	Arguments:
	  pattern - search pattern (required) (eg. ethereum.mainnet.usdt.0x1234.0x5677.whale)`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}

		pattern := args[0]

		// Prepare the WebSocket URL
		hurl := fmt.Sprintf("%s/v1/transfers/list/%s?page=%d&pageSize=%d", config.GetHost(), pattern, page, perPage)

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

		switch format {
		case "table":
			err = PrintTokenResponse(b)
		case "csv":
			err = RenderTransfersToCSV(b)
		case "json":
			fmt.Printf("%s", b)
		}

		if err != nil {
			log.Printf("failed to render response into table %s\n", err)
			return
		}
	},
}

func init() {
	ListCmd.Flags().IntVar(&page, "page", 0, "Page of returned results")
	ListCmd.Flags().IntVar(&perPage, "perPage", 20, "SizeBucket of page")
	ListCmd.Flags().StringVar(&format, "format", "table", "Output format (table,json,csv)")
}
