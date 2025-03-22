package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/heimdahl-xyz/heimdahl-cli/config"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type ContractParams struct {
	ProjectName     string  `json:"project_name"`
	Chain           string  `json:"chain"`
	Network         string  `json:"network"`
	ContractAddress string  `json:"contract_address"`
	ContractName    string  `json:"contract_name"`
	EventNames      *string `json:"event_names"`
	RawABI          *string `json:"raw_abi,omitempty"`
}

var (
	eventNames string
	rawABI     string
	rawABIFile string
)

var AddCmd = &cobra.Command{
	Use:   "add [address] [name]",
	Short: "Add a new contract",
	Long: `Add a new contract to the system.

Arguments:
  address - The contract address (required)
  name   - A user-defined name for the contract (required)`,
	Args: cobra.ExactArgs(2), // Expect exactly 2 arguments
	Run: func(cmd *cobra.Command, args []string) {
		contractAddress := args[0]
		contractName := args[1]

		params := ContractParams{
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

		req.Header.Set("Authorization", "Bearer "+config.GetApiKey())
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

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Contract already added %s", contractAddress)
			return
		}

		fmt.Printf("Successfully added contract %s", contractAddress)
	},
}

func init() {

	AddCmd.Flags().StringVarP(&rawABI, "abi", "r", "", "Raw ABI (optional)")
	AddCmd.Flags().StringVarP(&rawABIFile, "abi_file", "f", "", "Raw ABI file path (optional)")
}
