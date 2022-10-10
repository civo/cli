package network

import (
	"errors"

	"github.com/spf13/cobra"
)

// NetworkCmd manages Civo networks
var NetworkCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"networks", "net"},
	Short:   "Details of Civo networks",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

var networkSubnetCmd = &cobra.Command{
	Use:     "subnet",
	Aliases: []string{"subnets"},
	Short:   "Details of Civo network subnets",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cmd.Help()
		if err != nil {
			return err
		}
		return errors.New("a valid subcommand is required")
	},
}

func init() {
	NetworkCmd.AddCommand(networkListCmd)
	NetworkCmd.AddCommand(networkCreateCmd)
	NetworkCmd.AddCommand(networkUpdateCmd)
	NetworkCmd.AddCommand(networkRemoveCmd)
	NetworkCmd.AddCommand(networkSubnetCmd)

	networkSubnetCmd.AddCommand(networkSubnetListCmd)
	networkSubnetCmd.AddCommand(networkSubnetCreateCmd)
	networkSubnetCmd.AddCommand(networkSubnetRemoveCmd)
}
