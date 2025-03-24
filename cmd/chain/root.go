package chain

import (
	"github.com/spf13/cobra"
)

var ChainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Chain subcommands",
}

func init(){
	ChainCmd.AddCommand(ShowCmd)
	ChainCmd.AddCommand(ListCmd)
}
