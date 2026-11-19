package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/pkg/pluralize"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	deleteVolumes           bool
	keepFirewalls           bool
	keepKubeconfig          bool
	deleteKubeconfigContext bool // deprecated; kubeconfig context is removed by default
)

var kubernetesRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo kubernetes remove CLUSTER_NAME",
	Short:   "Remove a Kubernetes cluster",
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
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		targets := make([]clusterRemovalTarget, 0, len(args))
		if len(args) == 1 {
			kubernetesCluster, err := client.FindKubernetesCluster(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s Kubernetes cluster in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one Kubernetes cluster with that name in your account")
					os.Exit(1)
				}
				utility.Error("%s", err)
				os.Exit(1)
			}
			targets = append(targets, buildRemovalTarget(client, kubernetesCluster))
		} else {
			for _, v := range args {
				kubernetesCluster, err := client.FindKubernetesCluster(v)
				if err != nil {
					utility.Warning("skipping %q: %s", v, err)
					continue
				}
				targets = append(targets, buildRemovalTarget(client, kubernetesCluster))
			}
		}

		if len(targets) == 0 {
			utility.Error("no matching Kubernetes clusters found")
			os.Exit(1)
		}

		kubernetesNameList := make([]string, 0, len(targets))
		var keptVolumes []pendingVolume
		for _, t := range targets {
			kubernetesNameList = append(kubernetesNameList, t.Name)
			if t.VolumesError != "" {
				utility.Warning("could not list volumes for cluster %s before delete: %s. Attached volumes may be left behind with no leftover report", t.Name, t.VolumesError)
			}
			if !deleteVolumes && len(t.Volumes) > 0 {
				keptVolumes = append(keptVolumes, t.Volumes...)
			}
		}

		if len(keptVolumes) > 0 {
			utility.YellowConfirm("There are volume(s) attached to this cluster. They will be left behind unless you pass --delete-volumes:\n%s\n", utility.Green(strings.Join(volumeNames(keptVolumes), "\n")))
		}

		if !utility.UserConfirmedDeletion(fmt.Sprintf("Kubernetes %s", pluralize.Pluralize(len(targets), "cluster")), common.DefaultYes, strings.Join(kubernetesNameList, ", ")) {
			fmt.Println("Operation aborted.")
			return
		}

		removeKubeconfig := shouldRemoveKubeconfigContext(keepKubeconfig)
		if deleteKubeconfigContext && keepKubeconfig {
			removeKubeconfig = false
		} else if deleteKubeconfigContext {
			removeKubeconfig = true
		}

		for _, t := range targets {
			_, err = client.DeleteKubernetesCluster(t.ID)
			if err != nil {
				utility.Error("error deleting the kubernetes cluster: %s", err)
				os.Exit(1)
			}

			// Only strip the local context after the API accepted the delete, so a
			// failed delete leaves kubectl access intact.
			if removeKubeconfig {
				removeClusterKubeconfigContext(t.Name)
			}

			if deleteVolumes {
				var failed []pendingVolume
				for _, vol := range t.Volumes {
					if err := deleteClusterVolume(client, vol); err != nil {
						utility.Warning("could not delete volume %s (%s): %s", vol.Name, vol.ID, err)
						failed = append(failed, vol)
						continue
					}
					fmt.Printf("Deleted volume %s (%s)\n", utility.Green(vol.Name), vol.ID)
				}
				if len(failed) > 0 {
					fmt.Fprint(os.Stderr, leftBehindVolumesInfo(t.Name, failed))
				}
			} else if len(t.Volumes) > 0 {
				fmt.Fprint(os.Stderr, leftBehindVolumesInfo(t.Name, t.Volumes))
			}

			if shouldDeleteFirewall(keepFirewalls) && t.FirewallID != "" {
				if !t.FirewallExclusive {
					utility.Info("keeping firewall %s for cluster %s because it is shared with other resources; pass nothing extra, or delete it later with 'civo firewall remove %s'", t.FirewallID, t.Name, t.FirewallID)
				} else if _, err := client.DeleteFirewall(t.FirewallID); err != nil {
					// Cluster delete is async; the firewall may still be in use briefly.
					utility.Warning("could not delete firewall %s for cluster %s: %s. Delete it manually with 'civo firewall remove %s' once the cluster is gone, or pass --keep-firewalls to skip this step", t.FirewallID, t.Name, err, t.FirewallID)
				} else {
					fmt.Printf("Deleted firewall %s\n", utility.Green(t.FirewallID))
				}
			}
		}

		ow := utility.NewOutputWriter()

		for _, t := range targets {
			ow.StartLine()
			ow.AppendDataWithLabel("id", t.ID, "ID")
			ow.AppendDataWithLabel("name", t.Name, "Name")
		}

		switch common.OutputFormat {
		case "json":
			if len(targets) == 1 {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteMultipleObjectsJSON(common.PrettySet)
			}
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The Kubernetes %s (%s) %s been deleted\n",
				pluralize.Pluralize(len(targets), "cluster"),
				utility.Green(strings.Join(kubernetesNameList, ", ")),
				pluralize.Has(len(targets)),
			)
		}
	},
}

func buildRemovalTarget(client *civogo.Client, cluster *civogo.KubernetesCluster) clusterRemovalTarget {
	target := clusterRemovalTarget{
		ID:         cluster.ID,
		Name:       cluster.Name,
		FirewallID: cluster.FirewallID,
	}

	if cluster.FirewallID != "" {
		fw, err := client.FindFirewall(cluster.FirewallID)
		if err != nil {
			utility.Warning("could not inspect firewall %s for cluster %s: %s; it will not be auto-deleted", cluster.FirewallID, cluster.Name, err)
		} else {
			target.FirewallExclusive = firewallSafeToDelete(fw.ClusterCount, fw.InstanceCount, fw.LoadBalancerCount)
		}
	}

	vols, err := client.ListVolumesForCluster(cluster.ID)
	if err != nil {
		target.VolumesError = err.Error()
		return target
	}
	for _, v := range vols {
		target.Volumes = append(target.Volumes, pendingVolume{ID: v.ID, Name: v.Name})
	}
	return target
}

func removeClusterKubeconfigContext(clusterName string) {
	context := strings.ToLower(clusterName)
	var cmd *exec.Cmd
	if utility.CheckOS() == "windows" {
		cmd = exec.Command("powershell", "kubectl", "config", "delete-context", context)
	} else {
		cmd = exec.Command("kubectl", "config", "delete-context", context)
	}
	if err := cmd.Run(); err != nil {
		utility.Info("Kubeconfig context for cluster %s was not found (nothing to remove)", clusterName)
		return
	}
	fmt.Printf("Removed kubeconfig context for cluster %s\n", utility.Green(clusterName))
}

func deleteClusterVolume(client *civogo.Client, vol pendingVolume) error {
	// Best-effort detach first so a still-attached volume can be removed without
	// waiting for the async cluster delete to finish.
	_, _ = client.DetachVolume(vol.ID)
	_, err := client.DeleteVolume(vol.ID)
	return err
}
