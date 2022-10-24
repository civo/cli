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

		var network *civogo.NetworkResult

		nc := civogo.NetworkConfig{
			Label:       args[0],
			Region:      client.Region,
			IPv4Enabled: &v4enabled,
			IPv6Enabled: &v6enabled,
		}

		if cidrv4 != "" {
			nc.CIDRv4 = cidrv4
		}

		network, err = client.CreateNetwork(nc)
		if err != nil {
			utility.Error("Error creating the network: %s", err)
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
