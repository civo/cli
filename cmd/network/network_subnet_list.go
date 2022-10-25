package network

import (
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var networkSubnetListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: `civo network subnet ls <NETWORK_ID>"`,
	Args:    cobra.MinimumNArgs(1),
	Short:   "List all subnets within a network",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			return
		}

		subnets, err := client.ListSubnets(args[0])
		if err != nil {
			utility.Error("%s", err)
			return
		}

		ow := utility.NewOutputWriter()

		for _, subnet := range subnets {
			ow.StartLine()
			ow.AppendDataWithLabel("id", utility.TruncateID(subnet.ID), "ID")
			ow.AppendDataWithLabel("name", subnet.Name, "Name")
			ow.AppendDataWithLabel("subnet_size", subnet.SubnetSize, "Subnet Size")
			ow.AppendDataWithLabel("region", client.Region, "Region")
			ow.AppendDataWithLabel("status", subnet.Status, "Status")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("network_id", subnet.NetworkID, "Network ID")
			}
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
