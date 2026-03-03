package vpc

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var subnetCreateNetworkID string

var vpcSubnetCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo vpc subnet create SUBNET_NAME --network NETWORK_NAME",
	Short:   "Create a new VPC subnet",
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

		network, err := client.FindVPCNetwork(subnetCreateNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnet, err := client.CreateVPCSubnet(network.ID, civogo.SubnetConfig{
			Name: args[0],
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": subnet.ID, "name": subnet.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a VPC subnet called %s with ID %s\n", utility.Green(subnet.Name), utility.Green(subnet.ID))
		}
	},
}
