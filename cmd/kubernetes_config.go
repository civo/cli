package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var saveConfig, mergeConfig bool
var localPathConfig string

var kubernetesConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Example: "civo kubernetes config CLUSTER_NAME --save --merge",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Get a Kubernetes cluster's config",
	Long: `Show the Kubernetes config for a specified cluster.
If you wish to use a custom format, the available fields are:

	* KubeConfig`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		kube, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("Finding the Kubernetes cluster failed with %s", err)
			os.Exit(1)
		}

		if saveConfig {
			if !mergeConfig && strings.Contains(localPathConfig, ".kube") {
				if utility.UserConfirmedOverwrite("kubernetes config", defaultYes) == true {
					err := utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, mergeConfig, kube.Name)
					if err != nil {
						utility.Error("Saving the cluster config failed with %s", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("Operation aborted.")
					os.Exit(1)
				}
			} else {
				err := utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, mergeConfig, kube.Name)
				if err != nil {
					utility.Error("Saving the cluster config failed with %s", err)
					os.Exit(1)
				}
			}

		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"KubeConfig": kube.KubeConfig})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			if !saveConfig {
				fmt.Println(kube.KubeConfig)
			}

		}
	},
}
