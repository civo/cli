package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNodePoolListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all node pools in a Kubernetes cluster",
	Example: "civo kubernetes node-pool ls CLUSTER_NAME",
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

		for _, pool := range cluster.RequiredPools {
			ow = utility.NewOutputWriter()
			fmt.Println()
			ow.WriteHeader(fmt.Sprintf("Node Pool %s", pool.ID))
			ow.AppendDataWithLabel("ID", pool.ID, "Name")
			ow.AppendDataWithLabel("Size", pool.Size, "Size")
			ow.AppendDataWithLabel("Count", fmt.Sprintf("%d", pool.Count), "Count")
			labels, err := json.Marshal(pool.Labels)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			ow.AppendDataWithLabel("Labels", string(labels), "Labels")
			taints, err := json.Marshal(pool.Taints)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			ow.AppendDataWithLabel("Taints", string(taints), "Taints")
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
