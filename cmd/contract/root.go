package contract

import (
	"github.com/spf13/cobra"
)

var ContractCmd = &cobra.Command{
	Use:   "contract",
	Short: "Contract subcommands",
}

func init() {
	ContractCmd.AddCommand(AddCmd)
	ContractCmd.AddCommand(ShowCmd)
	ContractCmd.AddCommand(ListCmd)
}
