package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var targetNodesPoolSize string
var numTargetNodesPool int

var kubernetesNodePoolCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"create", "add"},
	Short:   "Add a node pool to Kubernetes cluster",
	Example: "civo kubernetes node-pool create CLUSTER_NAME [flags]",
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		newPool := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range kubernetesFindCluster.Pools {
			newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
		}

		poolID := uuid.NewString()
		newPool = append(newPool, civogo.KubernetesClusterPoolConfig{ID: poolID, Count: numTargetNodesPool, Size: targetNodesPoolSize})

		configKubernetes := &civogo.KubernetesClusterConfig{
			Pools: newPool,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The pool (%s) was added to the cluster (%s)\n", utility.Green(poolID[:6]), utility.Green(kubernetesCluster.Name))
		}
	},
}
