package cmd

import (
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kuberneteNodePoolList []utility.ObjecteList
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
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			kuberneteNodePoolList = append(kuberneteNodePoolList, utility.ObjecteList{ID: args[0], Name: args[1]})
		} else {
			for _, v := range args[1:] {
				kuberneteNodePoolList = append(kuberneteNodePoolList, utility.ObjecteList{ID: args[0], Name: v})
			}
		}

		kubernetesPoolNameList := []string{}
		for _, v := range kuberneteNodePoolList {
			kubernetesPoolNameList = append(kubernetesPoolNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("node %s", pluralize.Pluralize(len(kuberneteNodePoolList), "pool")), defaultYes, strings.Join(kubernetesPoolNameList, ", ")) {

			nodePool := []civogo.KubernetesClusterPoolConfig{}
			kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
			if err != nil {
				utility.Error("Kubernetes %s", err)
				os.Exit(1)
			}

			for _, v := range kubernetesFindCluster.Pools {
				nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
			}

			for _, kubeList := range kuberneteNodePoolList {
				nodePool = utility.RemoveNodePool(nodePool, kubeList.Name)
			}

			configKubernetes := &civogo.KubernetesClusterConfig{
				Pools: nodePool,
			}

			kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			ow := utility.NewOutputWriter()

			for _, v := range kuberneteNodePoolList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.Name, "ID")
			}

			switch outputFormat {
			case "json":
				if len(kuberneteNodePoolList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted from the cluster (%s)\n", fmt.Sprintf("node %s", pluralize.Pluralize(len(kuberneteNodePoolList), "pool")), utility.Green(strings.Join(kubernetesPoolNameList, ", ")), utility.Green(kubernetesCluster.Name))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
