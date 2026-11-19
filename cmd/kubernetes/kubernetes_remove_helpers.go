package kubernetes

import (
	"fmt"
	"strings"
)

// pendingVolume is a volume attached to a cluster at the time of deletion.
// IDs are captured before DeleteKubernetesCluster so we can still delete them
// afterwards (ListVolumesForCluster requires the cluster to still exist).
type pendingVolume struct {
	ID   string
	Name string
}

// clusterRemovalTarget is one cluster queued for deletion, with the dependent
// resources we may clean up after the API accepts the delete.
type clusterRemovalTarget struct {
	ID                string
	Name              string
	FirewallID        string
	FirewallExclusive bool // false when shared / unknown — never auto-delete then
	Volumes           []pendingVolume
	VolumesError      string // non-empty when inventory failed before delete
}

// shouldRemoveKubeconfigContext reports whether the cluster's kubeconfig
// context should be removed. Removal is the default; --keep-kubeconfig opts out.
func shouldRemoveKubeconfigContext(keepKubeconfig bool) bool {
	return !keepKubeconfig
}

// shouldDeleteFirewall reports whether the cluster firewall should be deleted.
// Deletion is the default; --keep-firewalls opts out.
func shouldDeleteFirewall(keepFirewalls bool) bool {
	return !keepFirewalls
}

// volumeNames returns display names for confirmation / leftover messages.
func volumeNames(volumes []pendingVolume) []string {
	names := make([]string, 0, len(volumes))
	for _, v := range volumes {
		if v.Name != "" {
			names = append(names, v.Name)
			continue
		}
		names = append(names, v.ID)
	}
	return names
}

// leftBehindVolumesInfo is the stderr hint when volumes are kept (or failed to
// delete) after a cluster delete.
func leftBehindVolumesInfo(clusterName string, volumes []pendingVolume) string {
	if len(volumes) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Volumes left behind for Kubernetes cluster %s:\n", clusterName)
	for _, v := range volumes {
		label := v.Name
		if label == "" {
			label = v.ID
		}
		fmt.Fprintf(&b, "- %s (%s)\n", label, v.ID)
	}
	b.WriteString("Use '--delete-volumes' next time to delete them automatically.\n")
	return b.String()
}

// firewallSafeToDelete reports whether a firewall looks exclusive to this
// cluster (not shared with other clusters, instances, or load balancers).
// Callers should evaluate this before requesting cluster deletion.
func firewallSafeToDelete(clusterCount, instanceCount, loadBalancerCount int) bool {
	return clusterCount <= 1 && instanceCount == 0 && loadBalancerCount == 0
}
