package swap

import (
	"github.com/spf13/cobra"
)

var SwapCmd = &cobra.Command{
	Use:   "swap",
	Short: "Swap tokens subcommands",
}

func init() {
	SwapCmd.AddCommand(ListCmd)
}
