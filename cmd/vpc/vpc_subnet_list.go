package vpc

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var subnetListNetworkID string

var vpcSubnetListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: "civo vpc subnet ls --network NETWORK_NAME",
	Short:   "List VPC subnets",
	Long: `List all subnets in a VPC network.
If you wish to use a custom format, the available fields are:

	* id
	* name
	* network_id
	* subnet_size
	* status`,
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

		network, err := client.FindVPCNetwork(subnetListNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnets, err := client.ListVPCSubnets(network.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, subnet := range subnets {
			ow.StartLine()
			ow.AppendDataWithLabel("id", subnet.ID, "ID")
			ow.AppendDataWithLabel("name", subnet.Name, "Name")
			ow.AppendDataWithLabel("network_id", subnet.NetworkID, "Network ID")
			ow.AppendDataWithLabel("subnet_size", subnet.SubnetSize, "Subnet Size")
			ow.AppendDataWithLabel("status", subnet.Status, "Status")
		}

		ow.FinishAndPrintOutput()
	},
}
