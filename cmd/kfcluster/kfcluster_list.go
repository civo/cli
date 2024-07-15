package kfcluster

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kfcListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo kfc ls`,
	Short:   "List all kubeflow clusters",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		kfclusters, err := client.ListKfClusters()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, kfc := range kfclusters.Items {
			ow.StartLine()
			ow.AppendDataWithLabel("id", kfc.ID, "ID")
			ow.AppendDataWithLabel("name", kfc.Name, "Name")
			ow.AppendDataWithLabel("size", kfc.Size, "Size")
			ow.AppendDataWithLabel("kubeflow_ready", kfc.KubeflowReady, "Kubeflow Ready")
			ow.AppendDataWithLabel("dashboard_url", kfc.DashboardURL, "Dashboard URL")
		}

		ow.FinishAndPrintOutput()
	},
}
