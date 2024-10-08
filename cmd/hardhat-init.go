package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/heimdahl-xyz/heimdahl-cli/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootFolder string
	networks   string
)

func createFileIfNotExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
		fmt.Println("Config created:", filename)
	} else if err != nil {
		return fmt.Errorf("Error checking file: %w", err)
	} else {
		fmt.Println("Config already exists:", filename)
	}
	return nil
}

var hardhatInitCmd = &cobra.Command{
	Use:   "hardhat-init",
	Short: "Init hardhat project configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// check if the file already exists
		var root string
		if rootFolder == "" {
			var err error
			root, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current working directory:", err)
				return
			}
		} else {
			root = rootFolder
		}

		err := createFileIfNotExists(root + "/heimdahl.json")
		if err != nil {
			fmt.Println("Error creating heimdahl.json file:", err)
			return
		}

		cf, err := os.ReadFile(root + "/package.json")
		if err != nil {
			fmt.Println("Error reading package.json file:", err)
			return
		}

		var m map[string]interface{}
		err = json.Unmarshal(cf, &m)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		metas, err := internal.GetContractMetas(rootFolder)
		if err != nil {
			log.Println("Error getting contract metas:", err)
			return
		}

		contractParams := make([]internal.ContractParams, 0)
		for _, met := range metas {
			args := make([]internal.SolidityArgument, 0)
			for _, inp := range met.ABI.Constructor.Inputs {
				args = append(args, internal.SolidityArgument{
					Name:  inp.Name,
					Type:  inp.Type.String(),
					Value: "-- PLEASE AMEND HERE ---",
				})
			}

			contractParams = append(contractParams, internal.ContractParams{
				ContractName: met.ContractName,
				Arguments: internal.ConstructorArguments{
					Inputs: args,
				},
			})

		}

		projectName := m["name"].(string)

		var networks = []internal.Network{
			{
				Chain:          "ethereum",
				Network:        "localnet",
				ChainID:        31337,
				PrivateKey:     "--PLEASE ADD YOUR PRIVATE KEY HERE--",
				ContractParams: contractParams,
			},
		}

		config := internal.Config{
			ProjectName: projectName,
			Networks:    networks,
		}

		b, _ := json.MarshalIndent(config, "", "  ")

		err = os.WriteFile(root+"/heimdahl.json", b, 0644)
		if err != nil {
			fmt.Println("Error writing to heimdahl.json file:", err)
			return
		}
	},
}

func init() {
	hardhatInitCmd.Flags().StringVarP(&rootFolder, "rootFolder", "r", ".", "Hardhat project root folder (optional)")
	_ = hardhatInitCmd.MarkFlagRequired("networks")
}
