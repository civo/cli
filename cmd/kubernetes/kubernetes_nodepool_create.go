package kubernetes

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var targetNodesPoolSize, nodePoolName string
var numTargetNodesPool int
var publicIpNodePool bool

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

		var poolID string
		if nodePoolName != "" {
			poolID = nodePoolName
		} else {
			poolID = uuid.NewString()
		}

		if len(poolID) > 63 {
			utility.Error("The pool name must be less than 64 characters")
			os.Exit(1)
		}

		poolConfig := civogo.KubernetesClusterPoolUpdateConfig{ID: poolID, Count: &numTargetNodesPool, Size: targetNodesPoolSize}
		if publicIpNodePool {
			if config.Current.RegionToFeatures != nil {
				if !config.Current.RegionToFeatures[client.Region].PublicIPNodePools {
					utility.Error("The region \"%s\" does not support \"Public IP Node Pools\" feature", client.Region)
					os.Exit(1)
				}
			}
			poolConfig.PublicIPNodePool = publicIpNodePool
		}

		kubernetesCluster, err := client.CreateKubernetesClusterPool(kubernetesFindCluster.ID, &civogo.KubernetesClusterPoolConfig{
			Region:           client.Region,
			ID:               poolID,
			Count:            numTargetNodesPool,
			Size:             targetNodesPoolSize,
			PublicIPNodePool: publicIpNodePool,
		})
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kubernetesCluster.ID, "name": kubernetesFindCluster.Name, "pool_id": poolID[:6]})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The pool (%s) was added to the cluster (%s)\n", utility.Green(poolID), utility.Green(kubernetesFindCluster.Name))
		}
	},
}
