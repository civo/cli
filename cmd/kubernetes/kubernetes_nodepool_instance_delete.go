package kubernetes

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var nodePoolID, instanceID string

var kubernetesNodePoolInstanceDeleteCmd = &cobra.Command{
	Use:     "instance-delete",
	Aliases: []string{"instance-rm", "instance-remove"},
	Short:   "Delete an instance from a node pool in a Kubernetes cluster",
	Example: "civo kubernetes node-pool instance-delete CLUSTER_NAME [flags]",
	Args:    cobra.MinimumNArgs(1),
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

		kubernetesClusterName := args[0]

		kubernetesFindCluster, err := client.FindKubernetesCluster(kubernetesClusterName)
		if err != nil {
			utility.Error("Kubernetes %s", err)
			os.Exit(1)
		}

		nodePool, err := client.FindKubernetesClusterPool(kubernetesFindCluster.ID, nodePoolID)
		if err != nil {
			utility.Error("Node pool %s", err)
			os.Exit(1)
		}
		if nodePool == nil {
			utility.Error("Node pool with ID %s not found in cluster %s", nodePoolID, kubernetesClusterName)
			os.Exit(1)
		}

		instance := findInstanceInPoolByID(nodePool.Instances, instanceID)
		if instance == nil {
			utility.Error("Instance with ID %s not found in node pool %s", instanceID, nodePoolID)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("instance %s", instanceID), common.DefaultYes, "") {

			_, err := client.DeleteKubernetesClusterPoolInstance(kubernetesFindCluster.ID, nodePoolID, instanceID)
			if err != nil {
				utility.Error("Error deleting instance %s from node pool %s in cluster %s: %s", instanceID, nodePoolID, kubernetesFindCluster.Name, err)
				os.Exit(1)
			}

			utility.Printf("Instance %s has been deleted from node pool %s in cluster %s\n", instanceID, nodePoolID, kubernetesFindCluster.Name)
		} else {
			utility.Println("Operation aborted.")
		}
	},
}

func findInstanceInPoolByID(instances []civogo.KubernetesInstance, instanceID string) *civogo.KubernetesInstance {
	for _, instance := range instances {
		if instance.ID == instanceID {
			return &instance
		}
	}
	return nil
}
