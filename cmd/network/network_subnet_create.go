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

var networkSubnetCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"create", "add"},
	Short:   "Create a new subnet",
	Example: "civo network subnet create <SUBNET-NAME> <NETWORK-ID>",
	Args:    cobra.MinimumNArgs(2),
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

		network, err := client.FindNetwork(args[1])
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		subnetConfig := civogo.SubnetConfig{
			Name: args[0],
		}

		subnet, err := client.CreateSubnet(network.ID, subnetConfig)
		if err != nil {
			utility.Error("Subnet %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": subnet.ID, "name": subnet.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The subnet (%s) was created in network with ID (%s)\n", utility.Green(subnet.Name), utility.Green(network.ID))
		}
	},
}
