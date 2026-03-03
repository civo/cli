package vpc

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var subnetShowNetworkID string

var vpcSubnetShowCmd = &cobra.Command{
	Use:     "show [SUBNET-NAME/SUBNET-ID]",
	Short:   "Show details of a specific VPC subnet",
	Aliases: []string{"get", "describe", "inspect"},
	Args:    cobra.ExactArgs(1),
	Example: "civo vpc subnet show SUBNET_NAME --network NETWORK_NAME",
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

		network, err := client.FindVPCNetwork(subnetShowNetworkID)
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnet, err := client.FindVPCSubnet(args[0], network.ID)
		if err != nil {
			utility.Error("Subnet %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", subnet.ID, "ID")
		ow.AppendDataWithLabel("name", subnet.Name, "Name")
		ow.AppendDataWithLabel("network_id", subnet.NetworkID, "Network ID")
		ow.AppendDataWithLabel("subnet_size", subnet.SubnetSize, "Subnet Size")
		ow.AppendDataWithLabel("status", subnet.Status, "Status")

		switch common.OutputFormat {
		case "json":
			ow.ToJSON(subnet, common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Println("VPC Subnet Details:")
			fmt.Printf("ID: %s\n", subnet.ID)
			fmt.Printf("Name: %s\n", subnet.Name)
			fmt.Printf("Network ID: %s\n", subnet.NetworkID)
			fmt.Printf("Subnet Size: %s\n", subnet.SubnetSize)
			fmt.Printf("Status: %s\n", subnet.Status)
		}
	},
}
