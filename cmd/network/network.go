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

	// Add VLAN-related flags
	networkCreateCmd.Flags().IntVar(&vlanID, "vlan-id", 0, "VLAN ID for the network")
	networkCreateCmd.Flags().StringVar(&vlanCIDRV4, "vlan-cidr-v4", "", "CIDR for VLAN IPv4")
	networkCreateCmd.Flags().StringVar(&vlanGatewayIPv4, "vlan-gateway-ip-v4", "", "Gateway IP for VLAN IPv4")
	networkCreateCmd.Flags().StringVar(&vlanHardwareAddr, "vlan-hardware-addr", "", "Hardware address for VLAN")
	networkCreateCmd.Flags().StringVar(&vlanAllocationStartV4, "vlan-allocation-pool-v4-start", "", "Start of the IPv4 allocation pool for VLAN")
	networkCreateCmd.Flags().StringVar(&vlanAllocationEndV4, "vlan-allocation-pool-v4-end", "", "End of the IPv4 allocation pool for VLAN")
}
