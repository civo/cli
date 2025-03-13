package kubernetes

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesUpdateKubeconfigCmd = &cobra.Command{
	Use:     "update-kubeconfig",
	Short:   "Update kubeconfig for the specified cluster",
	Example: "civo kubernetes update-kubeconfig CLUSTER_NAME",
	Args:    cobra.ExactArgs(1),
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

		cluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Finding cluster failed with %s", err)
			os.Exit(1)
		}

		err = utility.ObtainKubeConfig(localPathConfig, cluster.KubeConfig, true, true, cluster.Name)
		if err != nil {
			utility.Error("Updating kubeconfig failed with %s", err)
			os.Exit(1)
		}

		utility.Printf("Updated kubeconfig with cluster %s configuration\n", utility.Green(cluster.Name))
	},
}
