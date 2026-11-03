package kubernetes

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/pluralize"
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
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		requested := args[1:]
		for _, id := range requested {
			if len(id) < 6 {
				utility.Error("Please provide the node pool ID with at least 6 characters for %s", id)
				os.Exit(1)
			}
		}

		if !utility.UserConfirmedDeletion(fmt.Sprintf("node %s", pluralize.Pluralize(len(requested), "pool")), common.DefaultYes, strings.Join(requested, ", ")) {
			fmt.Println("Operation aborted.")
			return
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Kubernetes %s", err)
			os.Exit(1)
		}

		poolIDs, err := resolveNodePoolsToDelete(kubernetesFindCluster.Pools, requested)
		if err != nil {
			utility.Error("%s in cluster %s", err, kubernetesFindCluster.Name)
			os.Exit(1)
		}

		deleted, err := deleteNodePools(client, kubernetesFindCluster.ID, poolIDs)
		if err != nil {
			if len(deleted) > 0 {
				utility.Warning("Deleted node %s %s before the error", pluralize.Pluralize(len(deleted), "pool"), strings.Join(deleted, ", "))
			}
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		for _, id := range deleted {
			ow.StartLine()
			ow.AppendDataWithLabel("id", id, "ID")
		}

		switch common.OutputFormat {
		case "json":
			if len(deleted) == 1 {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteMultipleObjectsJSON(common.PrettySet)
			}
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The node %s (%s) %s been deleted from the cluster (%s)\n",
				pluralize.Pluralize(len(deleted), "pool"),
				utility.Green(strings.Join(deleted, ", ")),
				pluralize.Has(len(deleted)),
				utility.Green(kubernetesFindCluster.Name),
			)
		}
	},
}

// resolveNodePoolsToDelete maps each requested node pool ID prefix to the full
// ID of a pool in the cluster. It fails, without deleting anything, when a
// prefix matches no pool or when the deletion would leave the cluster with no
// pools (the API refuses to delete the last pool, so we would otherwise fail
// halfway through).
func resolveNodePoolsToDelete(pools []civogo.KubernetesPool, requested []string) ([]string, error) {
	remaining := make([]civogo.KubernetesClusterPoolConfig, 0, len(pools))
	for _, p := range pools {
		remaining = append(remaining, civogo.KubernetesClusterPoolConfig{ID: p.ID})
	}

	var ids, notFound []string
	for _, id := range requested {
		before := len(ids)
		remaining, ids = utility.RemoveNodePool(remaining, id, ids)
		if len(ids) == before {
			notFound = append(notFound, id)
		}
	}

	if len(notFound) > 0 {
		return nil, fmt.Errorf("no node pool found matching %s", strings.Join(notFound, ", "))
	}
	if len(remaining) == 0 {
		return nil, fmt.Errorf("cannot delete every node pool, a cluster needs at least one")
	}

	return ids, nil
}

type nodePoolDeleter interface {
	DeleteKubernetesClusterPool(id, poolID string) (*civogo.SimpleResponse, error)
}

// deleteNodePools deletes each pool through the dedicated pool endpoint, so the
// remaining pools keep their labels and taints. It stops at the first failure
// and returns the IDs deleted so far.
func deleteNodePools(client nodePoolDeleter, clusterID string, poolIDs []string) ([]string, error) {
	deleted := make([]string, 0, len(poolIDs))
	for _, id := range poolIDs {
		if _, err := client.DeleteKubernetesClusterPool(clusterID, id); err != nil {
			return deleted, fmt.Errorf("deleting node pool %s: %w", id, err)
		}
		deleted = append(deleted, id)
	}
	return deleted, nil
}
