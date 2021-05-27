package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"

	kubemartcmd "github.com/kubemart/kubemart-cli/cmd"
	kubemartutils "github.com/kubemart/kubemart-cli/pkg/utils"
)

var kubernetesAppAddCmd = &cobra.Command{
	Use:     "add",
	Example: "civo kubernetes application add NAME:PLAN --cluster CLUSTER_NAME\ncivo kubernetes application add NAME:\"LONG PLAN\" --cluster CLUSTER_NAME",
	Aliases: []string{"install"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Add the marketplace application to a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		isKubemart, err := isKubemartCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}

		if isKubemart {
			kubemartutils.DebugPrintf("Entering Kubemart branch\n")
			kubemartAdd(cmd, args)
		} else {
			legacyAdd(args)
		}
	},
}

func kubemartAdd(cmd *cobra.Command, args []string) {
	if kubernetesClusterApp == "" {
		utility.Error("Please provide --cluster or -c flag")
		os.Exit(1)
	}

	err := syncKubemartApps()
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}

	processedAppsAndPlanLabels, err := kubemartcmd.PreRunInstall(cmd, args)
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

	kubemartutils.DebugPrintf("Creating Kubemart App CRs\n")
	err = cs.RunInstall(processedAppsAndPlanLabels)
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}
}

func legacyAdd(args []string) {
	client, err := config.CivoAPIClient()
	if regionSet != "" {
		client.Region = regionSet
	}
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	kubernetesFindCluster, err := client.FindKubernetesCluster(kubernetesClusterApp)
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	appList, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	result := utility.RequestedSplit(appList, args[0])
	configKubernetes := &civogo.KubernetesClusterConfig{
		Applications: result,
	}

	kubeCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubeCluster.ID, "Name": kubeCluster.Name})

	switch outputFormat {
	case "json":
		ow.WriteSingleObjectJSON()
	case "custom":
		ow.WriteCustomOutput(outputFields)
	default:
		fmt.Printf("The application was installed in the Kubernetes cluster %s\n", utility.Green(kubeCluster.Name))
	}
}
