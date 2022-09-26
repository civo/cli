package kubernetes

import (
	"fmt"

	"github.com/civo/cli/common"
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

	* id
	* name
	* region
	* nodes
	* pools
	* status`,
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

		kubernetesClusters, err := client.ListKubernetesClusters()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, cluster := range kubernetesClusters.Items {
			ow.StartLine()

			ow.AppendDataWithLabel("id", cluster.ID, "ID")
			ow.AppendDataWithLabel("name", cluster.Name, "Name")
			ow.AppendDataWithLabel("region", client.Region, "Region")
			ow.AppendDataWithLabel("nodes", strconv.Itoa(len(cluster.Instances)), "Nodes")
			ow.AppendDataWithLabel("pools", strconv.Itoa(len(cluster.Pools)), "Pools")
			ow.AppendDataWithLabel("status", utility.ColorStatus(cluster.Status), "Status")

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("status", cluster.Status, "Status")
			}

			if config.Current.Clusters[cluster.ID] == false {
				ow.AppendDataWithLabel("name", cluster.Name+" *", "Name")
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
		fmt.Println()
		utility.Info("Cluster names marked with * are clusters which are over 1 year old. You might want to consider redownloading the config for these clusters. Ignore if already downloaded.")
	},
}
