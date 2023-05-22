package kubernetes

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNodePoolInstanceListCmd = &cobra.Command{
	Use:     "instance-ls",
	Aliases: []string{"instance-list", "instance-all"},
	Short:   "List all instances in a Kubernetes node pool",
	Example: "civo kubernetes node-pool instance-ls CLUSTER_NAME [flags]",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		cluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		if nodePoolID == "" {
			for _, pool := range cluster.Pools {
				for _, instance := range pool.Instances {

					ow.StartLine()
					ow.AppendDataWithLabel("ID", instance.ID, "ID")
					ow.AppendDataWithLabel("Hostname", instance.Hostname, "Hostname")
					ow.AppendDataWithLabel("Size", instance.Size, "Size")
					ow.AppendDataWithLabel("Status", instance.Status, "Status")
					ow.AppendDataWithLabel("Node Pool ID", pool.ID, "Node Pool ID")
				}
			}
		} else {
			pool, err := client.FindKubernetesClusterPool(cluster.ID, nodePoolID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			if pool == nil {
				utility.Error("Node pool with ID %s not found in cluster %s", nodePoolID, args[0])
				os.Exit(1)
			}
			for _, instance := range pool.Instances {
				ow.StartLine()
				ow.AppendDataWithLabel("ID", instance.ID, "ID")
				ow.AppendDataWithLabel("Hostname", instance.Hostname, "Hostname")
				ow.AppendDataWithLabel("Size", instance.Size, "Size")
				ow.AppendDataWithLabel("Status", instance.Status, "Status")
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
