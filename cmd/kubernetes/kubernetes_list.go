package kubernetes

import (
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
			ow.AppendDataWithLabel("cluster_type", cluster.ClusterType, "Cluster-Type")
			ow.AppendDataWithLabel("nodes", strconv.Itoa(len(cluster.Instances)), "Nodes")
			ow.AppendDataWithLabel("pools", strconv.Itoa(len(cluster.Pools)), "Pools")

			if cluster.Conditions != nil {
				conditions := ""
				for _, condition := range cluster.Conditions {
					if condition.Type == ControlPlaneReady {
						conditions += "Control Plane Accessible: " + string(condition.Status) + "\n"
					}
					if condition.Type == WorkerNodesReady {
						conditions += "All Workers Up: " + string(condition.Status) + "\n"
					}
					if condition.Type == ClusterVersionSync {
						conditions += "Cluster On Desired Version: " + string(condition.Status) + "\n"
					}
				}
				ow.AppendDataWithLabel("conditions", conditions, "Conditions")
			}

			if common.OutputFormat == "json" || common.OutputFormat == "custom" {
				ow.AppendDataWithLabel("status", cluster.Status, "Status")
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
