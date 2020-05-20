package cmd

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"networks", "net"},
	Short:   "Details of Civo networks",
}

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkUpdateCmd)
	networkCmd.AddCommand(networkRemoveCmd)
}
