package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type EventListenerParams struct {
	Network         string  `json:"network"`
	ContractAddress string  `json:"contract_address"`
	ContractName    string  `json:"contract_name"`
	EventNames      *string `json:"event_names"`
	RawABI          *string `json:"raw_abi,omitempty"`
}

var (
	network         string
	contractAddress string
	contractName    string
	eventNames      string
	rawABI          string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new event listener",
	Run: func(cmd *cobra.Command, args []string) {
		params := EventListenerParams{
			Network:         network,
			ContractAddress: contractAddress,
			ContractName:    contractName,
			EventNames:      &eventNames,
			RawABI:          &rawABI,
		}

		jsonData, err := json.Marshal(params)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		url := fmt.Sprintf("%s/api/v1/event-listeners", host) // Use the global host variable
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error making POST request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Event listener created successfully")
		} else {
			fmt.Printf("Failed to create event listener: %s\n", resp.Status)
		}
	},
}

func init() {
	createCmd.Flags().StringVarP(&network, "network", "n", "", "Blockchain network (required)")
	createCmd.Flags().StringVarP(&contractAddress, "address", "a", "", "Contract address (required)")
	createCmd.Flags().StringVarP(&contractName, "name", "c", "", "Contract name (required)")
	createCmd.Flags().StringVarP(&eventNames, "events", "e", "", "Comma-separated event names (optional)")
	createCmd.Flags().StringVarP(&rawABI, "abi", "r", "", "Raw ABI (optional)")

	createCmd.MarkFlagRequired("network")
	createCmd.MarkFlagRequired("address")
	createCmd.MarkFlagRequired("name")
}
