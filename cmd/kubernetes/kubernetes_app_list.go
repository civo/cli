package kubernetes

import (
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesAppListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all Kubernetes clusters applications",
	Long: `List all available Kubernetes clusters applications.
If you wish to use a custom format, the available fields are:

	* name
	* version
	* category
	* plans
	* dependencies`,
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

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

			ow.AppendDataWithLabel("name", kubeApp.Name, "Name")
			ow.AppendDataWithLabel("version", kubeApp.Version, "Version")
			ow.AppendDataWithLabel("category", kubeApp.Category, "Category")
			ow.AppendDataWithLabel("plans", strings.Join(plansApps, ", "), "Plans")
			ow.AppendDataWithLabel("dependencies", strings.Join(kubeApp.Dependencies, ", "), "Dependencies")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
