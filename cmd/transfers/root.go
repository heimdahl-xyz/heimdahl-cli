package transfers

import (
	"github.com/spf13/cobra"
)

var TransferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer tokens subcommands",
}

func init() {
	//TransferCmd.AddCommand(SubscribeCmd)
	TransferCmd.AddCommand(ListCmd)
	//EventCmd.AddCommand(ListCmd)
	//EventCmd.AddCommand(CrossListenCmd)
}
