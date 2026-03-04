package vpc

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var vpcNetworkUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"rename", "change", "modify"},
	Example: "civo vpc network update OLD_NAME NEW_NAME",
	Short:   "Rename a VPC network",
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

		oldNetwork, err := client.FindVPCNetwork(args[0])
		if err != nil {
			utility.Error("Network %s", err)
			os.Exit(1)
		}

		network, err := client.RenameVPCNetwork(args[1], oldNetwork.ID)
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
			fmt.Printf("Renamed the VPC network called %s with ID %s to %s\n", utility.Green(oldNetwork.Label), utility.Green(network.ID), utility.Green(network.Label))
		}
	},
}
