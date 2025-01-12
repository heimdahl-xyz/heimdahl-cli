package event

import (
	"github.com/spf13/cobra"
)

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Event subcommands",
}

func init() {
	EventCmd.AddCommand(ListenCmd)
	EventCmd.AddCommand(ListCmd)
	//EventCmd.AddCommand(CrossListenCmd)
}
