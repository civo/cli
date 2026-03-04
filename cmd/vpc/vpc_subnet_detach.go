package vpc

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var subnetDetachNetworkID string

var vpcSubnetDetachCmd = &cobra.Command{
	Use:     "detach [SUBNET-NAME/SUBNET-ID]",
	Aliases: []string{"disconnect"},
	Example: "civo vpc subnet detach SUBNET_NAME --network NETWORK_NAME",
	Short:   "Detach a VPC subnet from an instance",
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

		network, err := client.FindVPCNetwork(subnetDetachNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnet, err := client.FindVPCSubnet(args[0], network.ID)
		if err != nil {
			utility.Error("Subnet %s", err)
			os.Exit(1)
		}

		_, err = client.DetachVPCSubnetFromInstance(network.ID, subnet.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Detached VPC subnet %s from instance\n", utility.Green(subnet.Name))
		}
	},
}
