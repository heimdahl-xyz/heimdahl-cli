package subscription

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/url"
	"regexp"
	"strings"
)

type Topic string

type Subscription struct {
	StreamType string  `json:"stream_type"`
	Endpoint   string  `json:"endpoint"`
	Topics     []Topic `json:"topics"`
}

type TransferTopic struct {
	Chain   string `json:"chain"`
	Network string `json:"network"`
	Token   string `json:"token"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Size    string `json:"size"`
}

// Subscription struct to hold user input
type AddTransferSubscriptionParams struct {
	Identifier string
	Endpoint   string
	Chain      string
	Network    string
	Token      string
	From       string
	Wallet     string
	To         string
	SizeBucket string
}

func (t *TransferTopic) String() Topic {
	return Topic(fmt.Sprintf("%s.%s.%s.%s.%s.%s",
		t.Chain,
		t.Network,
		t.Token,
		t.From,
		t.To,
		t.Size))
}

// ValidateURL checks if the URL contains a valid scheme and either an IP address or hostname
func ValidateURL(input string) error {
	// Parse the URL using net/url package
	parsedURL, err := url.Parse(input)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Check if the URL has a valid scheme (http, https, etc.)
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a valid scheme (e.g., http:// or https://)")
	}

	// Check if the host is a valid IP or a valid hostname
	if !isValidIP(parsedURL.Host) && !isValidHostname(parsedURL.Host) {
		return fmt.Errorf("host must be a valid IP or hostname")
	}

	return nil
}

// isValidIP checks if the string is a valid IP address
func isValidIP(host string) bool {
	// Use net.ParseIP to check for valid IP addresses
	return net.ParseIP(host) != nil
}

// isValidHostname checks if the string is a valid hostname
func isValidHostname(host string) bool {
	// Remove port if present
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// Regular expression for valid hostname (simplified)
	hostnameRegex := `^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`
	match, _ := regexp.MatchString(hostnameRegex, host)
	return match
}

var AddTransferCmd = &cobra.Command{
	Use:   "add-transfer [address] [name]",
	Short: "Add a new transfer subscription",
	//Args:  cobra.ExactArgs(2), // Expect exactly 2 arguments
	Run: func(cmd *cobra.Command, args []string) {

		// Initialize the subscription struct
		var subscription AddTransferSubscriptionParams
		subscription.SizeBucket = "all"
		// Prompt for StreamType
		prompt := promptui.Prompt{
			Label: "Identifier (e.g., unified_token_transfers)",
		}
		prompt.Validate = func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("subscription id is required")
			}
			return nil
		}

		id, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		subscription.Identifier = strings.Trim(id, " ")

		// Prompt for Endpoint URL
		prompt = promptui.Prompt{
			Label: "Endpoint URL",
		}

		prompt.Validate = func(s string) error {
			if err = ValidateURL(s); err != nil {
				return fmt.Errorf("invalid endpoint %s %s", s, err)
			}
			return nil
		}

		endpoint, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		subscription.Endpoint = endpoint

		// Prompt for Chain
		prompt = promptui.Prompt{
			Label: "Chain (e.g., ethereum, base, binance, polygon, solana, tron)",
		}

		chain, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		subscription.Chain = strings.Trim(chain, " ")

		// Prompt for Network
		prompt = promptui.Prompt{
			Label: "Network (e.g., mainnet)",
		}
		network, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		subscription.Network = strings.Trim(network, " ")

		// Prompt for Token
		prompt = promptui.Prompt{
			Label: "Token (e.g., USDC, USDT, DAI)",
		}
		token, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}
		subscription.Token = strings.ToLower(strings.Trim(token, " "))

		// Prompt for the type of subscription (either from/to or wallet)

		sprompt := promptui.Select{
			Label: "Choose address input type",
			Items: []string{"From/To Addresses", "Wallet Address"},
		}
		_, result, err := sprompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		// Conditional prompts based on user's selection
		if result == "From/To Addresses" {
			// Prompt for From address
			prompt = promptui.Prompt{
				Label: "From address (eg. 0xaxxxxx default: all)",
			}
			from, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}
			subscription.From = from

			// Prompt for To address
			prompt = promptui.Prompt{
				Label: "To address (eg. 0xaxxxxx default: all)",
			}
			to, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}
			subscription.To = to
		} else if result == "Wallet Address" {
			// Prompt for Wallet address
			prompt = promptui.Prompt{
				Label: "Wallet address (eg. 0xaxxxxx default: all)",
			}
			wallet, err := prompt.Run()
			if err != nil {
				log.Fatalf("Prompt failed %v\n", err)
			}
			subscription.Wallet = wallet
		}

		fmt.Printf("\nSubscription Created: %+v\n", subscription)

		//contractAddress := args[0]
		//contractName := args[1]
		//
		//params := ContractParams{
		//	Chain:           chain,
		//	Network:         network,
		//	ContractAddress: contractAddress,
		//	ContractName:    contractName,
		//}
		//
		//if rawABIFile != "" {
		//	abb, err := os.ReadFile(rawABIFile)
		//	if err != nil {
		//		fmt.Println("Error reading ABI file:", err)
		//		return
		//	}
		//
		//	_, err = abi.JSON(bytes.NewReader(abb))
		//	if err != nil {
		//		fmt.Println("Error parsing ABI:", err)
		//		return
		//	}
		//	rawABI = string(abb)
		//}
		//
		//if rawABI != "" {
		//	params.RawABI = &rawABI
		//}
		//
		//if eventNames != "" {
		//	params.EventNames = &eventNames
		//}
		//
		//jsonData, err := json.Marshal(params)
		//if err != nil {
		//	fmt.Println("Error marshalling JSON:", err)
		//	return
		//}
		//
		//url := fmt.Sprintf("%s/v1/contracts", config.GetHost()) // Use the global host variable
		//
		//req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
		//if err != nil {
		//	fmt.Printf("Error creating post request %s:", err)
		//	return
		//}
		//
		//req.Header.Set("Authorization", "Bearer "+config.GetApiKey())
		//req.Header.Set("Content-Type", "application/json")
		//
		//resp, err := http.DefaultClient.Do(req)
		//if err != nil {
		//	fmt.Println("Error making POST request:", err)
		//	return
		//}
		//defer resp.Body.Close()
		//
		//if !(resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK) {
		//	fmt.Println("Failed to create event listener", resp.Status)
		//	return
		//}
		//
		//if resp.StatusCode == http.StatusOK {
		//	fmt.Printf("Contract already added %s", contractAddress)
		//	return
		//}
		//
		//fmt.Printf("Successfully added contract %s", contractAddress)
	},
}

func init() {

	//AddTransferCmd.Flags().StringVarP(&rawABI, "abi", "r", "", "Raw ABI (optional)")
	//AddTransferCmd.Flags().StringVarP(&rawABIFile, "abi_file", "f", "", "Raw ABI file path (optional)")
}
