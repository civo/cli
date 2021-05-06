package cmd

import (
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var kubernetesListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo kubernetes ls -o custom -f "ID: Name"`,
	Short:   "List all Kubernetes clusters",
	Long: `List all Kubernetes clusters.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Region
	* Nodes
	* Pools
	* Status`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()

		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesClusters, err := client.ListKubernetesClusters()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, cluster := range kubernetesClusters.Items {
			ow.StartLine()

			ow.AppendData("ID", cluster.ID)
			ow.AppendData("Name", cluster.Name)
			ow.AppendData("Region", client.Region)
			ow.AppendData("Nodes", strconv.Itoa(len(cluster.Instances)))
			ow.AppendData("Pools", strconv.Itoa(len(cluster.Pools)))
			ow.AppendData("Status", utility.ColorStatus(cluster.Status))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("Status", cluster.Status)
			}

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
