package kubernetes

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var numTargetNodesPoolScale int

var kubernetesNodePoolScaleCmd = &cobra.Command{
	Use:     "scale",
	Short:   "Scale a node pool in a Kubernetes cluster",
	Example: "civo kubernetes node-pool scale CLUSTER_NAME NODEPOOL_ID [flags]",
	Args:    cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		nodePoolID := args[1]
		if len(nodePoolID) < 6 {
			utility.Error("Please provide the node pool ID with at least 6 characters for %s", nodePoolID)
			os.Exit(1)
		}

		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Kubernetes %s", err)
			os.Exit(1)
		}

		nodePoolFound := false
		for _, pool := range kubernetesFindCluster.Pools {
			if strings.Contains(pool.ID, nodePoolID) {
				nodePoolID = pool.ID
				nodePoolFound = true
			}
		}

		if !nodePoolFound {
			utility.Error("Unable to find %q node pool inside %q cluster", nodePoolID, kubernetesFindCluster.Name)
			os.Exit(1)
		}

		nodePool := []civogo.KubernetesClusterPoolConfig{}
		for _, v := range kubernetesFindCluster.Pools {
			nodePool = append(nodePool, civogo.KubernetesClusterPoolConfig{ID: v.ID, Count: v.Count, Size: v.Size})
		}

		nodePool = utility.UpdateNodePool(nodePool, nodePoolID, numTargetNodesPoolScale)

		configKubernetes := &civogo.KubernetesClusterConfig{
			Pools: nodePool,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kubernetesCluster.ID, "name": kubernetesCluster.Name, "pool_id": nodePoolID})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The pool (%s) was scaled to (%s) in the cluster (%s)\n", utility.Green(nodePoolID), utility.Green(strconv.Itoa(numTargetNodesPoolScale)), utility.Green(kubernetesCluster.Name))
		}
	},
}
