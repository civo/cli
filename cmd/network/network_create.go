package network

import (
	"fmt"
	"github.com/civo/civogo"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	cidrV4        string
	nameserversV4 []string

	// VLAN-related variables
	vlanID                int
	vlanCIDRV4            string
	vlanGatewayIPv4       string
	vlanHardwareAddr      string
	vlanAllocationStartV4 string
	vlanAllocationEndV4   string
)

var networkCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo network create NAME",
	Short:   "Create a new network",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		networkConfig := civogo.NetworkConfig{
			Label:         args[0],
			CIDRv4:        cidrV4,
			NameserversV4: nameserversV4,
		}

		vlanConnectConfig := &civogo.VLANConnectConfig{
			VlanID:                vlanID,
			HardwareAddr:          vlanHardwareAddr,
			CIDRv4:                vlanCIDRV4,
			GatewayIPv4:           vlanGatewayIPv4,
			AllocationPoolV4Start: vlanAllocationStartV4,
			AllocationPoolV4End:   vlanAllocationEndV4,
		}

		// Conditionally add VLAN configuration
		// Assuming VLAN ID of 0 is invalid and indicates no VLAN config provided
		if vlanID != 0 {
			networkConfig.VLanConfig = vlanConnectConfig
		}

		network, err := client.CreateNetwork(networkConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": network.ID, "label": network.Label})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a network called %s with ID %s\n", utility.Green(network.Label), utility.Green(network.ID))
		}
	},
}
