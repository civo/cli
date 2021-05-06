package cmd

import (
	"fmt"
	"os"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesNodePoolDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"delete", "rm"},
	Short:   "Delete a node pool from Kubernetes cluster",
	Example: "civo kubernetes node-pool delete CLUSTER_NAME NODEPOOL_ID [flags]",
	Args:    cobra.MinimumNArgs(2),
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
		if utility.UserConfirmedDeletion(fmt.Sprintf("node %s", pluralize.Pluralize(len(instanceList), "pool")), defaultYes, args[1]) {
			kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			nodePool := []civogo.KubernetesClusterPoolConfig{}
			for _, v := range kubernetesFindCluster.Pools {
				nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
			}
			nodePool = utility.RemoveNodePool(nodePool, args[1])
			configKubernetes := &civogo.KubernetesClusterConfig{
				Pools: nodePool,
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
				fmt.Printf("The pool (%s) was delete from the cluster (%s)\n", utility.Green(args[1]), utility.Green(kubernetesCluster.Name))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
