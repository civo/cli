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

func init() {
	NetworkCmd.AddCommand(networkListCmd)
	NetworkCmd.AddCommand(networkCreateCmd)
	NetworkCmd.AddCommand(networkUpdateCmd)
	NetworkCmd.AddCommand(networkRemoveCmd)

	networkCreateCmd.Flags().StringVarP(&cidrV4, "cidr-v4", "", "", "Custom IPv4 CIDR")
	networkCreateCmd.Flags().StringSliceVarP(&nameserversV4, "nameservers-v4", "", nil, "Custom list of IPv4 nameservers (up to three, comma-separated)")
}
