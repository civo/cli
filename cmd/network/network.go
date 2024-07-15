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
	NetworkCmd.AddCommand(networkShowCmd)
	NetworkCmd.AddCommand(networkConnectCmd)

	networkCreateCmd.Flags().StringVarP(&cidrV4, "cidr-v4", "", "", "Custom IPv4 CIDR")
	networkCreateCmd.Flags().StringSliceVarP(&nameserversV4, "nameservers-v4", "", nil, "Custom list of IPv4 nameservers (up to three, comma-separated)")
	networkCreateCmd.Flags().BoolVarP(&createDefaultFirewall, "create-default-firewall", "", false, "Create a default firewall for the network")

	// Flag binding for networkConnectCmd
	networkConnectCmd.Flags().IntVar(&vlanID, "vlan-id", 0, "VLAN ID to connect")
	networkConnectCmd.Flags().StringVar(&vlanCIDRV4, "cidr-v4", "", "CIDR v4 of the VLAN")
	networkConnectCmd.Flags().StringVar(&vlanGatewayIPv4, "gateway-ipv4", "", "Gateway IPv4 address for the VLAN")
	networkConnectCmd.Flags().StringVar(&vlanPhysicalInterface, "physical-interface", "eth0", "Physical interface for the VLAN connection")
	networkConnectCmd.Flags().StringVar(&vlanAllocationStartV4, "allocation-pool-v4-start", "", "Start of the IPv4 allocation pool for the VLAN")
	networkConnectCmd.Flags().StringVar(&vlanAllocationEndV4, "allocation-pool-v4-end", "", "End of the IPv4 allocation pool for the VLAN")
}
