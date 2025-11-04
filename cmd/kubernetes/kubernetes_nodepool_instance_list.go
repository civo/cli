package kubernetes

import (
	"fmt"
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
				fmt.Println()
				ow.WriteHeader(fmt.Sprintf("Node Pool %s", pool.ID))
				owPool := utility.NewOutputWriter()

				// Print all instances in this node pool
				for _, instance := range pool.Instances {
					owPool.StartLine()
					owPool.AppendDataWithLabel("ID", instance.ID, "ID")
					owPool.AppendDataWithLabel("Hostname", instance.Hostname, "Hostname")
					owPool.AppendDataWithLabel("Size", instance.Size, "Size")
					owPool.AppendDataWithLabel("Status", instance.Status, "Status")
					owPool.AppendDataWithLabel("Node Pool", pool.ID, "Node Pool")
				}
				owPool.WriteTable()
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
			ow.WriteTable()
		}
		if common.OutputFormat == "json" {
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		}

		if common.OutputFormat == "custom" {
			ow.WriteCustomOutput(common.OutputFields)
		}
	},
}
