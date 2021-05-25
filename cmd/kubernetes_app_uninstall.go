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

var kubernetesAppUninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Example: "civo kubernetes application uninstall NAME --cluster CLUSTER_NAME",
	Aliases: []string{"remove"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Uninstall a marketplace application from a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		isKubemart, err := isKubemartCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}

		if isKubemart {
			kubemartutils.DebugPrintf("Entering Kubemart branch\n")
			kubemartUninstall(args)
		} else {
			fmt.Println(legacyMarketplaceWarning)
		}
	},
}

func kubemartUninstall(args []string) {
	if kubernetesClusterApp == "" {
		utility.Error("Please provide --cluster or -c flag")
		os.Exit(1)
	}

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

	kubemartutils.DebugPrintf("Uninstall Kubemart apps\n")
	err = cs.RunUninstall(args)
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}
}
