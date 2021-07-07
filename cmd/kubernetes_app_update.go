package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var appsUpdateCmdExample = []string{
	"civo kubernetes application update NAME --cluster CLUSTER_NAME",
	"civo kubernetes application update NAME1,NAME2 --cluster CLUSTER_NAME",
}

var kubernetesAppUpdateCmd = &cobra.Command{
	Use:     "update",
	Example: strings.Join(appsUpdateCmdExample, "\n"),
	Args:    cobra.MinimumNArgs(1),
	Short:   "Update the marketplace application in a Kubernetes cluster",
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

		userInput := args[0]
		userInputSplitted := strings.Split(userInput, ",")
		for _, input := range userInputSplitted {
			appName := ""
			for _, installedApplication := range kubernetesCluster.InstalledApplications {
				if strings.EqualFold(installedApplication.Name, input) || strings.EqualFold(installedApplication.Title, input) {
					appName = installedApplication.Name
					break
				}
			}

			if appName == "" {
				utility.Error("%s", fmt.Errorf("app with name %s not found", input))
				os.Exit(1)
			}

			_, err = client.UpdateKubernetesApp(kubernetesCluster.ID, appName)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			fmt.Printf("%s is now scheduled for update process\n", input)
		}
	},
}
