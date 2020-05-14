package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var kubernetesListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo kubernetes ls -o custom -f "ID: Name"`,
	Short:   "List all kubernetes clusters",
	Long: `List all kubernetes clusters.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Node
	* Size
	* Status`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		kubernetesCluster, err := client.ListKubernetesClusters()
		if err != nil {
			utility.Error("Unable to list kubernetes cluster %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, kubernetes := range kubernetesCluster.Items {
			ow.StartLine()

			ow.AppendData("ID", kubernetes.ID)
			ow.AppendData("Name", kubernetes.Name)
			ow.AppendData("Node", strconv.Itoa(kubernetes.NumTargetNode))
			ow.AppendData("Size", kubernetes.TargetNodeSize)
			ow.AppendData("Status", fmt.Sprintf("%s", utility.ColorStatus(kubernetes.Status)))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("Status", kubernetes.Status)
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
