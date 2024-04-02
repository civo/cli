package network

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
	"os"
)

var networkShowCmd = &cobra.Command{
	Use:     "show [NETWORK-NAME/NETWORK-ID]",
	Short:   "Show details of a specific Civo network, including VLAN information if available",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		networkID := args[0]

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("failed to create Civo API client: %s", err)
			os.Exit(1)
		}

		network, err := client.GetNetwork(networkID)
		if err != nil {
			utility.Error("Failed to retrieve network: %s", err)
			os.Exit(1)
		}

		fmt.Println("Network Details:")
		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendData("ID", network.ID)
		ow.AppendData("Name", network.Name)
		ow.AppendData("Default", utility.BoolToYesNo(network.Default))
		ow.AppendData("CIDR", network.CIDR)
		ow.AppendData("Status", network.Status)
		ow.AppendData("IPv4 Enabled", utility.BoolToYesNo(network.IPv4Enabled))
		ow.AppendData("IPv6 Enabled", utility.BoolToYesNo(network.IPv6Enabled))
		ow.WriteTable()

		// Conditional VLAN Details
		if network.VLAN.VlanID != 0 {
			fmt.Println("\nVLAN Details:")
			ow = utility.NewOutputWriter() // Reset for a new section
			ow.StartLine()
			ow.AppendData("VLAN ID", fmt.Sprintf("%d", network.VLAN.VlanID))
			ow.AppendData("Hardware Address", network.VLAN.HardwareAddr)
			ow.AppendData("CIDRv4", network.VLAN.CIDRv4)
			ow.AppendData("Gateway IPv4", network.VLAN.GatewayIPv4)
			ow.AppendData("Allocation Pool IPv4 Start", network.VLAN.AllocationPoolV4Start)
			ow.AppendData("Allocation Pool IPv4 End", network.VLAN.AllocationPoolV4End)
			ow.WriteTable()
		} else {
			fmt.Println("No VLAN Configuration")
		}

		// Nameserver Details
		if len(network.NameserversV4) > 0 || len(network.NameserversV6) > 0 {
			fmt.Println("\nNameserver Details:")
			ow = utility.NewOutputWriter()
			ow.StartLine()
			if len(network.NameserversV4) > 0 {
				ow.AppendData("Nameservers IPv4", utility.SliceToString(network.NameserversV4))
			}
			if len(network.NameserversV6) > 0 {
				ow.AppendData("Nameservers IPv6", utility.SliceToString(network.NameserversV6))
			}
			ow.WriteTable()
		}

	},
}
