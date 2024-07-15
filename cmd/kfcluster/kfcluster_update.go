package kfcluster

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var updatedName string
var kfcUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"modify", "change"},
	Short:   "Update a kubeflow cluster",
	Example: "civo kfc update OLD_NAME --name NEW_NAME",
	Args:    cobra.MinimumNArgs(1),
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

		findKfCluster, err := client.FindKfCluster(args[0])
		if err != nil {
			utility.Error("KfCluster %s", err)
			os.Exit(1)
		}

		updatedKFC, err := client.UpdateKfCluster(findKfCluster.ID, &civogo.UpdateKfClusterReq{
			Name:   updatedName,
			Region: client.Region,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": updatedKFC.ID, "name": updatedKFC.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The KfCluster %s was updated\n", utility.Green(findKfCluster.Name))
			os.Exit(0)
		}
	},
}
