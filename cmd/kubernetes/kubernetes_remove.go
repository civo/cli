package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var (
	deleteVolumes  bool // Flag to delete dependent volumes
	keepFirewalls  bool // Flag to keep dependent firewalls
	keepKubeconfig bool // Flag to keep kubeconfig
)

var kuberneteList []utility.ObjecteList // List to store the Kubernetes clusters to be deleted

// Command definition for removing a Kubernetes cluster
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

		// Create a client to interact with the Civo API
		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		// Find and store the Kubernetes clusters to be deleted
		for _, clusterName := range args {
			kubernetesCluster, err := client.FindKubernetesCluster(clusterName)
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s Kubernetes cluster in your account", utility.Red(clusterName))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one Kubernetes cluster with that name in your account")
					os.Exit(1)
				}
			}
			kuberneteList = append(kuberneteList, utility.ObjecteList{ID: kubernetesCluster.ID, Name: kubernetesCluster.Name})
		}

		// Collect names of the clusters to be deleted for confirmation message
		kubernetesNameList := []string{}
		for _, v := range kuberneteList {
			kubernetesNameList = append(kubernetesNameList, v.Name)
		}

		// Confirm the deletion with the user
		if utility.UserConfirmedDeletion(fmt.Sprintf("Kubernetes %s", pluralize.Pluralize(len(kuberneteList), "cluster")), common.DefaultYes, strings.Join(kubernetesNameList, ", ")) {
			for _, v := range kuberneteList {

				// Delete the Kubernetes cluster
				_, err := client.DeleteKubernetesCluster(v.ID)
				if err != nil {
					utility.Error("error deleting the kubernetes cluster: %s", err)
					os.Exit(1)
				}

				/* Poll for the deletion status
				This is required because, if we try to delete other things like firewalls before
				the cluster is completely deleted, we will encounter errors like "database_firewall_inuse_by_cluster".
				*/
				for {
					_, err := client.FindKubernetesCluster(v.Name)
					if err != nil {
						if errors.Is(err, civogo.ZeroMatchesError) {
							break // Cluster is deleted
						}
					}
				}

				// Delete volumes if --delete-volumes flag is set
				if deleteVolumes {
					volumes, err := client.ListVolumesForCluster(v.ID)
					if err != nil {
						if !errors.Is(err, civogo.ZeroMatchesError) {
							utility.Error("Error listing volumes for cluster %s: %s", v.Name, err)
						}
					} else {
						for _, volume := range volumes {
							_, err := client.DeleteVolume(volume.ID)
							if err != nil {
								utility.Error("Error deleting volume %s: %s", volume.Name, err)
							}
							fmt.Printf("%s volume deleted", volume.ID)
						}
					}
				}

				// Output volumes left behind if deleteVolumes flag is not set
				if !deleteVolumes {
					volumes, err := client.ListVolumesForCluster(v.ID)
					if err != nil {
						if !errors.Is(err, civogo.ZeroMatchesError) {
							utility.Error("Error listing volumes for cluster %s: %s", v.Name, err)
						}
					} else if len(volumes) > 0 {
						fmt.Fprintf(os.Stderr, "Volumes left behind for Kubernetes cluster %s:\n", v.Name)
						for _, volume := range volumes {
							fmt.Fprintf(os.Stderr, "- %s\n", volume.ID)
						}
						fmt.Fprintf(os.Stderr, "Consider using '--delete-volumes' flag next time to delete them automatically.\n")
					}
				}

				// Delete firewalls if --keep-firewalls flag is not set
				if !keepFirewalls {
					firewall, err := client.FindFirewall(v.Name)
					if err != nil {
						if errors.Is(err, civogo.MultipleMatchesError) {
							utility.Error("Error deleting the firewall: %v. Please delete the firewall manually by visiting: https://dashboard.civo.com/firewalls", err)
						} else {
							utility.Error("Error finding the firewall: %v", err)
						}
					} else if firewall.ClusterCount == 0 {
						_, err := client.DeleteFirewall(firewall.ID)
						if err != nil {
							utility.Error("Error deleting firewall %s: %v", firewall.Name, err)
						} else {
							fmt.Fprintf(os.Stderr, "Firewall %s deleted. If you want to keep the firewall in the future, please use '--keep-firewalls' flag.\n", firewall.ID)
						}
					}
				}

				// Delete kubeconfig if --keep-kubeconfig flag is not set
				if !keepKubeconfig {
					kubeconfigPath := fmt.Sprintf("%s/.kube/config", os.Getenv("HOME"))
					err := os.Remove(kubeconfigPath)
					if err != nil && !os.IsNotExist(err) {
						utility.Error("Error deleting kubeconfig: %s", err)
					} else {
						fmt.Fprintf(os.Stderr, "Kubeconfig file deleted. If you want to keep the kubeconfig in the future, please use '--keep-kubeconfig' flag.\n")
					}
				}
			}

			// Output the result of the deletion
			ow := utility.NewOutputWriter()
			for _, v := range kuberneteList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			// Format the output based on the selected format
			switch common.OutputFormat {
			case "json":
				if len(kuberneteList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The Kubernetes %s (%s) has been deleted\n", pluralize.Pluralize(len(kuberneteList), "cluster"), utility.Green(strings.Join(kubernetesNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
