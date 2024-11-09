package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
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

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new event listener contract",
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

		url := fmt.Sprintf("%s/v1/contracts", config.GetHost()) // Use the global host variable

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating post request %s:", err)
			return
		}

		req.Header.Set("X-API-Key", config.GetApiKey())
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
	AddCmd.Flags().StringVarP(&chain, "chain", "c", "", "Blockchain network (eg. ethereum, required)")
	AddCmd.Flags().StringVarP(&network, "network", "n", "mainnet", "Blockchain network (eg. mainnet, required)")
	AddCmd.Flags().StringVarP(&contractAddress, "address", "a", "", "Contract address (eg, 0xdAC17F958D2ee523a2206206994597C13D831ec7 required)")
	AddCmd.Flags().StringVarP(&contractName, "name", "N", "", "Contract name (eg \"USDC Token\", required)")
	AddCmd.Flags().StringVarP(&eventNames, "events", "e", "", "Comma-separated event names (optional)")
	AddCmd.Flags().StringVarP(&rawABI, "abi", "r", "", "Raw ABI (optional)")
	AddCmd.Flags().StringVarP(&rawABIFile, "abi_file", "f", "", "Raw ABI file path (optional)")

	AddCmd.MarkFlagRequired("chain")
	AddCmd.MarkFlagRequired("network")
	AddCmd.MarkFlagRequired("address")
	AddCmd.MarkFlagRequired("name")
}
