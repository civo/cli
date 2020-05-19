package cmd

import (
	_ "errors"
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	_ "github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	_ "strconv"
)

var saveConfig, mergeConfig bool
var localPathConfig string

var kubernetesConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Example: "civo kubernetes config CLUSTER_NAME --save --merge",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Get kubernetes clusters config",
	Long: `Show current kubernetes config.
If you wish to use a custom format, the available fields are:

	* KubeConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		kube, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Unable to find the kubernetes cluster %s", err)
			os.Exit(1)
		}

		if saveConfig {
			err := utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, mergeConfig)
			if err != nil {
				utility.Error("Unable to save the cluster config %s", err)
				os.Exit(1)
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"KubeConfig": kube.KubeConfig})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Println("The configuration was save")
		}
	},
}
