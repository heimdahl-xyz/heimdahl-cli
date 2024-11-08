package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type EventListenerParams struct {
	ProjectName     string  `json:"project_name"`
	Chain           string  `json:"chain"`
	Network         string  `json:"network"`
	ContractAddress string  `json:"contract_address"`
	ContractName    string  `json:"contract_name"`
	EventNames      *string `json:"event_names"`
	RawABI          *string `json:"raw_abi,omitempty"`
}

var (
	chain           string
	network         string
	contractAddress string
	contractName    string
	eventNames      string
	rawABI          string
	rawABIFile      string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new event listener",
	Run: func(cmd *cobra.Command, args []string) {
		params := EventListenerParams{
			Chain:           chain,
			Network:         network,
			ContractAddress: contractAddress,
			ContractName:    contractName,
		}

		if rawABIFile != "" {
			abb, err := os.ReadFile(rawABIFile)
			if err != nil {
				fmt.Println("Error reading ABI file:", err)
				return
			}

			_, err = abi.JSON(bytes.NewReader(abb))
			if err != nil {
				fmt.Println("Error parsing ABI:", err)
				return
			}
			rawABI = string(abb)
		}

		if rawABI != "" {
			params.RawABI = &rawABI
		}

		if eventNames != "" {
			params.EventNames = &eventNames
		}

		jsonData, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		url := fmt.Sprintf("%s/v1/event-listeners", getHost()) // Use the global host variable

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating post request %s:", err)
			return
		}

		req.Header.Set("X-API-Key", getApiKey())
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
			fmt.Println("Failed to create event listener", resp.Status)
			return
		}
		fmt.Printf("Successfully created event listener for contract %s", contractAddress)

	},
}

func init() {
	createCmd.Flags().StringVarP(&chain, "chain", "c", "", "Blockchain network (eg. ethereum, required)")
	createCmd.Flags().StringVarP(&network, "network", "n", "mainnet", "Blockchain network (eg. mainnet, required)")
	createCmd.Flags().StringVarP(&contractAddress, "address", "a", "", "Contract address (eg, 0xdAC17F958D2ee523a2206206994597C13D831ec7 required)")
	createCmd.Flags().StringVarP(&contractName, "name", "N", "", "Contract name (eg \"USDC Token\", required)")
	createCmd.Flags().StringVarP(&eventNames, "events", "e", "", "Comma-separated event names (optional)")
	createCmd.Flags().StringVarP(&rawABI, "abi", "r", "", "Raw ABI (optional)")
	createCmd.Flags().StringVarP(&rawABIFile, "abi_file", "f", "", "Raw ABI file path (optional)")

	createCmd.MarkFlagRequired("chain")
	createCmd.MarkFlagRequired("network")
	createCmd.MarkFlagRequired("address")
	createCmd.MarkFlagRequired("name")
}
