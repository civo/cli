package kubernetes

import (
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var firewall string

var kubernetesUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Update the configuration of your Civo Kubernetes cluster",
	Example: "civo kubernetes update CLUSTER_NAME --firewall FIREWALL_ID/NAME",
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

		if firewall == "" {
			utility.Error("You must specify either the firewall ID or the firewall name")
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			FirewallID: firewall,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": kubernetesCluster.ID, "name": kubernetesFindCluster.Name, "firewall_id": kubernetesCluster.FirewallID})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			utility.Printf("The firewall for the Kubernetes cluster %s has been updated\n", utility.Green(kubernetesFindCluster.Name))
		}
	},
}
