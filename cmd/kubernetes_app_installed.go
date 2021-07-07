package cmd

import (
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesAppInstalledCmd = &cobra.Command{
	Use:     "installed",
	Short:   "List installed Kubernetes applications in cluster",
	Example: "civo kubernetes application installed --cluster CLUSTER_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kubernetesCluster, err := client.FindKubernetesCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}
		kubeApps := kubernetesCluster.InstalledApplications

		ow := utility.NewOutputWriter()
		for _, kubeApp := range kubeApps {
			ow.StartLine()

			ow.AppendDataWithLabel("name", kubeApp.Name, "Name")
			ow.AppendDataWithLabel("version", kubeApp.Version, "Version")
			ow.AppendDataWithLabel("category", kubeApp.Category, "Category")
			ow.AppendDataWithLabel("plan", kubeApp.Plan, "Plan")
			ow.AppendDataWithLabel("dependencies", strings.Join(kubeApp.Dependencies, ", "), "Dependencies")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
