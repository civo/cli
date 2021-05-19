package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"

	kubemartcmd "github.com/kubemart/kubemart-cli/cmd"
	kubemartutils "github.com/kubemart/kubemart-cli/pkg/utils"
)

var kubernetesAppAddCmd = &cobra.Command{
	Use:     "add",
	Example: "civo kubernetes application add NAME:PLAN --cluster CLUSTER_NAME",
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
			kubemartAdd(args)
		} else {
			legacyAdd(args)
		}
	},
}

func kubemartAdd(args []string) {
	if kubernetesClusterApp == "" {
		utility.Error("Please provide --cluster or -c flag")
		os.Exit(1)
	}

	appNameWithPlan := args[0]
	splitted := strings.Split(appNameWithPlan, ":")
	appName := splitted[0]
	kubemartutils.DebugPrintf("App name to install: %s\n", appName)

	plan := 0
	if len(splitted) > 1 {
		plan = kubemartutils.ExtractPlanIntFromPlanStr(splitted[1])
	}

	kubemartutils.DebugPrintf("Calling Kubemart PreRunInstall with args appName: %s plan: %d\n", appName, plan)
	err := kubemartcmd.PreRunInstall(&appName, &plan)
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

	kubemartutils.DebugPrintf("Creating Kubemart App CR\n")
	err = cs.RunInstall(&appName, &plan)
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
