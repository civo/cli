package cmd

import (
	"fmt"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/gookit/color"
	kubemartcmd "github.com/kubemart/kubemart-cli/cmd"
	kubemartutils "github.com/kubemart/kubemart-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var kubernetesAppShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo kubernetes applications show APP_NAME CLUSTER_NAME"`,
	Args:    cobra.MinimumNArgs(2),
	Short:   "Shows the details of an application installed in the cluster",
	Long:    `Shows the details of an application installed in the cluster`,
	// ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	if len(args) == 0 {
	// 		return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
	// 	}
	// 	return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	// },
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		isKubemart, err := isKubemartCluster(args[1])
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}

		if isKubemart {
			kubemartutils.DebugPrintf("Entering Kubemart branch\n")
			kubemartShow(args)
		} else {
			legacyShow(args)
		}
	},
}

func kubemartShow(args []string) {
	appName := args[0]
	if appName == "" {
		utility.Error("Please provide app name")
		os.Exit(1)
	}

	clusterName := args[1]
	if clusterName == "" {
		utility.Error("Please provide cluster name")
		os.Exit(1)
	}

	err := syncKubemartApps()
	if err != nil {
		utility.Error(err.Error())
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
	kubernetesCluster, err := client.FindKubernetesCluster(clusterName)
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

	kubemartutils.DebugPrintf("Finding Kubemart app\n")
	_, err = cs.GetApp(appName)
	if err != nil {
		utility.Error("Sorry the app %s was not found in the cluster %s", appName, clusterName)
		os.Exit(1)
	}

	appPostInstall, err := kubemartutils.GetPostInstallMarkdown(appName)
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(appPostInstall)
}

func legacyShow(args []string) {
	client, err := config.CivoAPIClient()
	if regionSet != "" {
		client.Region = regionSet
	}
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	kubernetesCluster, err := client.FindKubernetesCluster(args[1])
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	foundAPP := false

	for _, app := range kubernetesCluster.InstalledApplications {
		if strings.EqualFold(app.Name, args[0]) {
			foundAPP = true
			result := markdown.Render(app.PostInstall, 80, 0)
			printPostInstall := color.S256()
			fmt.Println()
			printPostInstall.Println(string(result))
		}
	}

	if !foundAPP {
		utility.Error("Sorry the app %s was not found in the cluster %s", args[0], args[1])
		os.Exit(1)
	}
}
