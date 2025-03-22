package event

import (
	"github.com/spf13/cobra"
)

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Event subcommands",
}

func init() {
	EventCmd.AddCommand(SubscribeCmd)
	EventCmd.AddCommand(ListCmd)
}
