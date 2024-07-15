package network

import (
	"fmt"
	"os"

	"github.com/civo/civogo"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	networkCidrV4        string
	networkNameserversV4 []string

	// VLAN-related variables
	vlanID                int
	vlanCIDRV4            string
	vlanGatewayIPv4       string
	vlanPhysicalInterface string
	vlanAllocationStartV4 string
	vlanAllocationEndV4   string
)

var networkConnectCmd = &cobra.Command{
	Use:     "connect",
	Aliases: []string{"new", "add"},
	Example: "civo network connect VLAN_NAME",
	Short:   "Attach an existing VLAN to a network",
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
			CIDRv4:        networkCidrV4,
			NameserversV4: networkNameserversV4,
		}

		vlanConnectConfig := &civogo.VLANConnectConfig{
			VlanID:                vlanID,
			PhysicalInterface:     vlanPhysicalInterface,
			CIDRv4:                vlanCIDRV4,
			GatewayIPv4:           vlanGatewayIPv4,
			AllocationPoolV4Start: vlanAllocationStartV4,
			AllocationPoolV4End:   vlanAllocationEndV4,
		}

		networkConfig.VLanConfig = vlanConnectConfig

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
