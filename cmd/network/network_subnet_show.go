package network

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var networkSubnetShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "info"},
	Example: `civo network subnet show SUBNET_NAME NETWORK_ID`,
	Short:   "Prints information about a subnet",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		subnet, err := client.FindSubnet(args[0], args[1])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		ow.StartLine()
		fmt.Println()
		ow.AppendDataWithLabel("id", subnet.ID, "ID")
		ow.AppendDataWithLabel("name", subnet.Name, "Name")
		ow.AppendDataWithLabel("size", subnet.SubnetSize, "Subnet Size")
		ow.AppendDataWithLabel("network_id", subnet.NetworkID, "Network ID")
		ow.AppendDataWithLabel("region", client.Region, "Region")
		ow.AppendDataWithLabel("status", subnet.Status, "Status")

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
