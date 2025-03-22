package transfers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"io"
	"log"
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
	Timestamp    int64  `json:"timestamp"`
	FromAddress  string `json:"from_address"`
	ToAddress    string `json:"to_address"`
	Amount       int64  `json:"amount"`
	TokenAddress string `json:"token_address"`
	Symbol       string `json:"symbol"`
	Chain        string `json:"chain"`
	Network      string `json:"network"`
	TxHash       string `json:"tx_hash"`
	Decimals     int    `json:"decimals"`
	Position     int64  `json:"position"`
}

// TokenData represents the structure of the JSON data
type TokenData struct {
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

// formatAddress shortens an Ethereum address for display
func formatAddress(address string) string {
	if len(address) < 12 {
		return address
	}
	return address[:6] + "..." + address[len(address)-4:]
}

// formatAmount converts token amount based on decimals
func formatAmount(amount int64, decimals int) string {
	if decimals == 0 {
		return fmt.Sprintf("%d", amount)
	}

	divisor := int64(1)
	for i := 0; i < decimals; i++ {
		divisor *= 10
	}

	whole := amount / divisor
	fraction := amount % divisor

	// Format with appropriate trailing zeros
	fractionStr := fmt.Sprintf("%0*d", decimals, fraction)
	// Trim trailing zeros
	fractionStr = strings.TrimRight(fractionStr, "0")
	if fractionStr == "" {
		return fmt.Sprintf("%d", whole)
	}
	return fmt.Sprintf("%d.%s", whole, fractionStr)
}

// formatTimestamp converts Unix timestamp to human-readable time
func formatTimestamp(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// RenderTransfersTable renders the token transfers as a table
func RenderTransfersTable(jsonData []byte) error {
	var tokenData TokenData
	err := json.Unmarshal(jsonData, &tokenData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Define column widths
	cols := []struct {
		title string
		width int
	}{
		{"Time", 19},
		{"From", 42},
		{"To", 42},
		{"Amount", 12},
		{"Symbol", 6},
		{"Chain", 10},
		{"TX Hash", 66},
	}

	// Calculate total width for the table
	totalWidth := 1 // Starting with 1 for the left border
	for _, col := range cols {
		totalWidth += col.width + 3 // Width + 3 for padding and borders
	}

	// Print table header
	printLine(totalWidth)
	fmt.Print("|")
	for _, col := range cols {
		fmt.Printf(" %-*s |", col.width, col.title)
	}
	fmt.Println()
	printLine(totalWidth)

	// Print table rows
	for _, transfer := range tokenData.Transfers {
		fmt.Print("|")

		// Time
		fmt.Printf(" %-*s |", cols[0].width, formatTimestamp(transfer.Timestamp))

		// From
		fmt.Printf(" %s |", transfer.FromAddress)

		// To
		fmt.Printf(" %s |", transfer.ToAddress)

		// Amount
		fmt.Printf(" %-*s |", cols[3].width, formatAmount(transfer.Amount, transfer.Decimals))

		// Symbol
		fmt.Printf(" %-*s |", cols[4].width, transfer.Symbol)

		// Chain
		chainNetwork := transfer.Chain
		if transfer.Network != "" && transfer.Network != "mainnet" {
			chainNetwork += "." + transfer.Network
		}
		fmt.Printf(" %-*s |", cols[5].width, chainNetwork)

		// TX Hash
		fmt.Printf(" %s |", transfer.TxHash)

		fmt.Println()
	}

	// Print footer
	printLine(totalWidth)

	// Print metadata
	fmt.Println("\nMetadata:")
	fmt.Printf("Timestamp: %s\n", formatTimestamp(tokenData.Meta.Timestamp))
	fmt.Printf("Total transfers: %d\n", tokenData.Meta.Total)
	fmt.Printf("Page: %d of %d\n", tokenData.Meta.Page+1, (tokenData.Meta.Total+tokenData.Meta.PerPage-1)/tokenData.Meta.PerPage)
	fmt.Printf("Chains: %s\n", strings.Join(tokenData.Meta.Chains, ", "))
	fmt.Printf("Tokens: %s\n", strings.Join(tokenData.Meta.Tokens, ", "))

	return nil
}

// RenderTransfersToCSV exports token transfers as CSV to stdout
func RenderTransfersToCSV(jsonData []byte) error {
	var tokenData TokenData
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
			strconv.FormatInt(transfer.Amount, 10),
			formatAmount(transfer.Amount, transfer.Decimals),
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
			err = RenderTransfersTable(b)
		case "csv":
			err = RenderTransfersToCSV(b)
		case "json":
			fmt.Printf("%s", b)
		}

		if err != nil {
			log.Println("failed to render response into table %s", err)
			return
		}
	},
}

func init() {
	ListCmd.Flags().IntVar(&page, "page", 0, "Page of returned results")
	ListCmd.Flags().IntVar(&perPage, "perPage", 20, "Size of page")
	ListCmd.Flags().StringVar(&format, "format", "table", "Output format (table,json,csv)")
}
