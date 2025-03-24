package subscription

import (
	"github.com/spf13/cobra"
)

var SubscriptionCmd = &cobra.Command{
	Use:   "subscription",
	Short: "Subscription subcommands (Experimental)",
}

func init() {
	SubscriptionCmd.AddCommand(AddTransferCmd)
	SubscriptionCmd.AddCommand(ShowCmd)
	SubscriptionCmd.AddCommand(ListCmd)
}
