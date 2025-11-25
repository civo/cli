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

var kubernetesList []utility.Resource
var deleteKubeconfigContext bool

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
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

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
			}
			kubernetesList = append(kubernetesList, utility.Resource{ID: kubernetesCluster.ID, Name: kubernetesCluster.Name})
		} else {
			for _, v := range args {
				kubernetesCluster, err := client.FindKubernetesCluster(v)
				if err == nil {
					kubernetesList = append(kubernetesList, utility.Resource{ID: kubernetesCluster.ID, Name: kubernetesCluster.Name})
				}
			}
		}

		kubernetesNameList := []string{}
		for _, v := range kubernetesList {
			kubernetesNameList = append(kubernetesNameList, v.Name)
		}

		var volsNameList []string
		for _, v := range kubernetesList {
			vols, err := client.ListVolumesForCluster(v.ID)
			if err != nil {
				utility.Error("error getting the list of dangling volumes: %s", err)
				os.Exit(1)
			}
			for _, v := range vols {
				volsNameList = append(volsNameList, v.Name)
			}
			if vols != nil {
				utility.YellowConfirm("There are volume(s) attached to this cluster. Consider deleting or detaching these before deleting the cluster:\n%s\n", utility.Green(strings.Join(volsNameList, "\n")))
			}
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("Kubernetes %s", pluralize.Pluralize(len(kubernetesList), "cluster")), common.DefaultYes, strings.Join(kubernetesNameList, ", ")) {

			for _, v := range kubernetesList {
				if deleteKubeconfigContext {
					// Try to remove the kubeconfig context before deleting the cluster
					cmd := exec.Command("kubectl", "config", "delete-context", strings.ToLower(v.Name))
					if err := cmd.Run(); err != nil {
						fmt.Printf("Note: Kubeconfig context for cluster %s was not found\n", utility.Green(v.Name))
					} else {
						fmt.Printf("Successfully removed kubeconfig context for cluster %s\n", utility.Green(v.Name))
					}
				}
				_, err = client.DeleteKubernetesCluster(v.ID)
				if err != nil {
					utility.Error("error deleting the kubernetes cluster: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range kubernetesList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch common.OutputFormat {
			case "json":
				if len(kubernetesList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The Kubernetes %s (%s) %s been deleted\n",
					pluralize.Pluralize(len(kubernetesList), "cluster"),
					utility.Green(strings.Join(kubernetesNameList, ", ")),
					pluralize.Has(len(kubernetesList)),
				)
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
