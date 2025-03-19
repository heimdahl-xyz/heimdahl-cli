package fungibles

import (
	"github.com/spf13/cobra"
)

var FungibleCmd = &cobra.Command{
	Use:   "fungible",
	Short: "Fungible tokens subcommands",
}

func init() {
	//FungibleCmd.AddCommand(SubscribeCmd)
	FungibleCmd.AddCommand(ListCmd)
	//EventCmd.AddCommand(ListCmd)
	//EventCmd.AddCommand(CrossListenCmd)
}
