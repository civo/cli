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

var (
	subnetAttachNetworkID    string
	subnetAttachResourceID   string
	subnetAttachResourceType string
)

var vpcSubnetAttachCmd = &cobra.Command{
	Use:     "attach [SUBNET-NAME/SUBNET-ID]",
	Aliases: []string{"connect"},
	Example: "civo vpc subnet attach SUBNET_NAME --network NETWORK_NAME --resource-id RESOURCE_ID --resource-type instance",
	Short:   "Attach a VPC subnet to an instance",
	Args:    cobra.ExactArgs(1),
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

		network, err := client.FindVPCNetwork(subnetAttachNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnet, err := client.FindVPCSubnet(args[0], network.ID)
		if err != nil {
			utility.Error("Subnet %s", err)
			os.Exit(1)
		}

		route, err := client.AttachVPCSubnetToInstance(network.ID, subnet.ID, &civogo.CreateRoute{
			ResourceID:   subnetAttachResourceID,
			ResourceType: subnetAttachResourceType,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": route.ID, "subnet_id": subnet.ID, "resource_id": route.ResourceID})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Attached VPC subnet %s to resource %s\n", utility.Green(subnet.Name), utility.Green(subnetAttachResourceID))
		}
	},
}
