package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"

	kubemartcmd "github.com/kubemart/kubemart-cli/cmd"
	kubemartutils "github.com/kubemart/kubemart-cli/pkg/utils"
)

var kubernetesAppUpdateCmd = &cobra.Command{
	Use:     "update",
	Example: "civo kubernetes application update NAME --cluster CLUSTER_NAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Update a marketplace application in Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		isKubemart, err := isKubemartCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}

		if isKubemart {
			kubemartutils.DebugPrintf("Entering Kubemart branch\n")
			kubemartUpdate(args)
		} else {
			fmt.Println(legacyMarketplaceWarning)
		}
	},
}

func kubemartUpdate(args []string) {
	if kubernetesClusterApp == "" {
		utility.Error("Please provide --cluster or -c flag")
		os.Exit(1)
	}

	appName := args[0]
	kubemartutils.DebugPrintf("App name to update: %s\n", &appName)

	kubemartutils.DebugPrintf("Creating Civo API client\n")
	client, err := config.CivoAPIClient()
	if regionSet != "" {
		client.Region = regionSet
	}
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	kubemartutils.DebugPrintf("Finding Civo Kubernetes cluster\n")
	kubernetesCluster, err := client.FindKubernetesCluster(kubernetesClusterApp)
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	kubemartutils.DebugPrintf("Finding kubeconfig\n")
	kubeconfigStr := kubernetesCluster.KubeConfig

	kubemartutils.DebugPrintf("Creating Kubemart client from kubeconfig\n")
	cs, err := kubemartcmd.NewClientFromKubeConfigString(kubeconfigStr)
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}

	kubemartutils.DebugPrintf("Update a Kubemart app\n")
	err = cs.RunUpdate(&appName)
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}
}
