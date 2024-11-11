package stream

import (
	"github.com/spf13/cobra"
)

var StreamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream subcommands",
}

func init() {
	StreamCmd.AddCommand(ListenCmd)
	StreamCmd.AddCommand(ReplayCmd)
}
