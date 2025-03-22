package network

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
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

		// Display Core Network Details
		utility.Println("Network Details:")
		utility.Printf("ID: %s\n", network.ID)
		utility.Printf("Name: %s\n", network.Name)
		utility.Printf("Default: %s\n", utility.BoolToYesNo(network.Default))
		utility.Printf("CIDR: %s\n", network.CIDR)
		utility.Printf("Status: %s\n", network.Status)
		utility.Printf("IPv4 Enabled: %s\n", utility.BoolToYesNo(network.IPv4Enabled))
		utility.Printf("IPv6 Enabled: %s\n", utility.BoolToYesNo(network.IPv6Enabled))

		// Conditional VLAN Details
		if network.VlanID != 0 {
			utility.Println("\nVLAN Details:")
			utility.Printf("VLAN ID: %d\n", network.VlanID)
			utility.Printf("Hardware Address: %s\n", network.PhysicalInterface)
			utility.Printf("Gateway IPv4: %s\n", network.GatewayIPv4)
			utility.Printf("Allocation Pool IPv4 Start: %s\n", network.AllocationPoolV4Start)
			utility.Printf("Allocation Pool IPv4 End: %s\n", network.AllocationPoolV4End)
		} else {
			utility.Println("\nNo VLAN Configuration")
		}

		// Nameserver Details
		if len(network.NameserversV4) > 0 || len(network.NameserversV6) > 0 {
			utility.Println("\nNameserver Details:")
			if len(network.NameserversV4) > 0 {
				utility.Printf("Nameservers IPv4: %s\n", utility.SliceToString(network.NameserversV4))
			}
			if len(network.NameserversV6) > 0 {
				utility.Printf("Nameservers IPv6: %s\n", utility.SliceToString(network.NameserversV6))
			}
		}
	},
}
