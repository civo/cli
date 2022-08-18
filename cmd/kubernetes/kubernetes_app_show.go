package kubernetes

import (
	"fmt"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var kubernetesAppShowCmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"get", "inspect"},
	Example: `civo kubernetes applications show APP_NAME CLUSTER_NAME"`,
	Args:    cobra.MinimumNArgs(2),
	Short:   "Shows the details of an application installed in the cluster.",
	Long:    `Shows the details of an application installed in the cluster`,
	// ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	if len(args) == 0 {
	// 		return getAllKubernetesList(), cobra.ShellCompDirectiveNoFileComp
	// 	}
	// 	return getKubernetesList(toComplete), cobra.ShellCompDirectiveNoFileComp
	// },
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
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
			utility.Warning("Sorry the app %s was not found in the cluster %s", args[0], args[1])
			os.Exit(1)
		}

	},
}
