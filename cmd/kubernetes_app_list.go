package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	kubemartcmd "github.com/kubemart/kubemart-cli/cmd"
	"github.com/spf13/cobra"
	sortmap "github.com/tg/gosortmap"
)

var kubernetesAppListCmd = &cobra.Command{
	Use:     "ls",
	Example: "civo kubernetes application ls --cluster CLUSTER_NAME",
	Aliases: []string{"list", "all"},
	Short:   "List all Kubernetes clusters applications",
	Long: `List all available Kubernetes clusters applications.
If you wish to use a custom format, the available fields are:

	* Name
	* Version
	* Category
	* Plans
	* Dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		isKubemart, err := isKubemartCluster(kubernetesClusterApp)
		if err != nil {
			utility.Error(err.Error())
			os.Exit(1)
		}

		if isKubemart {
			kubemartList()
		} else {
			legacyList()
		}
	},
}

func kubemartList() {
	err := syncKubemartApps()
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}

	manifests, err := kubemartcmd.GetAppManifestsMap()
	if err != nil {
		utility.Error(err.Error())
		os.Exit(1)
	}

	ow := utility.NewOutputWriter()
	for _, m := range sortmap.ByKey(manifests) {
		ow.StartLine()
		mapKey := fmt.Sprintf("%v", m.Key)

		// All plans
		var plansApps []string
		for _, plan := range manifests[mapKey].Plans {
			plansApps = append(plansApps, plan.Label)
		}

		ow.AppendData("Name", mapKey)
		ow.AppendData("Version", manifests[mapKey].Version)
		ow.AppendData("Category", manifests[mapKey].Category)
		ow.AppendData("Plans", strings.Join(plansApps, ", "))
		ow.AppendData("Dependencies", strings.Join(manifests[mapKey].Dependencies, ", "))
	}

	switch outputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON()
	case "custom":
		ow.WriteCustomOutput(outputFields)
	default:
		ow.WriteTable()
	}
}

func legacyList() {
	client, err := config.CivoAPIClient()
	if err != nil {
		utility.Error("Creating the connection to Civo's API failed with %s", err)
		os.Exit(1)
	}

	kubeApps, err := client.ListKubernetesMarketplaceApplications()
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}

	ow := utility.NewOutputWriter()
	for _, kubeApp := range kubeApps {
		ow.StartLine()

		// All plans
		var plansApps []string
		for _, plan := range kubeApp.Plans {
			plansApps = append(plansApps, plan.Label)
		}

		ow.AppendData("Name", kubeApp.Name)
		ow.AppendData("Version", kubeApp.Version)
		ow.AppendData("Category", kubeApp.Category)
		ow.AppendData("Plans", strings.Join(plansApps, ", "))
		ow.AppendData("Dependencies", strings.Join(kubeApp.Dependencies, ", "))
	}

	switch outputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON()
	case "custom":
		ow.WriteCustomOutput(outputFields)
	default:
		ow.WriteTable()
	}
}
