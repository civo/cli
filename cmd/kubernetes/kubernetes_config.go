package kubernetes

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesConfigCmdExample = `* To merge and save:
    civo kubernetes config CLUSTER_NAME --save

* To merge and save, and switch to the new context:
	civo kubernetes config CLUSTER_NAME --save --switch

* To overwrite and save:
    civo kubernetes config CLUSTER_NAME --save --overwrite

Notes:
* By default, when --save is specified, we will merge your kubeconfig (unless --overwrite is specified).
* To auto-switch to new kubeconfig, --switch is required. Without it, your active context will remain unchanged.
* When --overwrite is specified, --switch is not required. Your context will be updated automatically.
`

var saveConfig, switchConfig, overwriteConfig bool
var localPathConfig string

var kubernetesConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Example: kubernetesConfigCmdExample,
	Args:    cobra.MinimumNArgs(1),
	Short:   "Get a Kubernetes cluster's config",
	Long: `Show the Kubernetes config for a specified cluster.
If you wish to use a custom format, the available fields are:

	* kubeconfig`,
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

		if switchConfig && overwriteConfig {
			utility.Info("--switch is not required when --overwrite is present")
		}

		kube, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if !kube.Ready {
			utility.Error("The cluster isn't ready yet, so the KUBECONFIG isn't available.")
			os.Exit(1)
		}

		if os.Getenv("KUBECONFIG") != "" {
			localPathConfig = os.Getenv("KUBECONFIG")
		}

		if saveConfig {
			if overwriteConfig {
				// overwrite and save
				if localPathConfig == os.Getenv("KUBECONFIG") {
					utility.Warning("Your $KUBECONFIG is set to %s, you are about to overwrite that config with your current kubeconfig.", localPathConfig)
				}
				if utility.UserConfirmedOverwrite("kubernetes config", common.DefaultYes) {
					err := utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, false, switchConfig, kube.Name)
					if err != nil {
						utility.Error("Saving the cluster config failed with %s", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("Operation aborted.")
					os.Exit(1)
				}
			} else {
				// merge and save
				err := utility.ObtainKubeConfig(localPathConfig, kube.KubeConfig, true, switchConfig, kube.Name)
				if err != nil {
					utility.Error("Saving the cluster config failed with %s", err)
					os.Exit(1)
				}
			}

		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"kubeconfig": kube.KubeConfig})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if !saveConfig {
				fmt.Println(kube.KubeConfig)
			}

		}
	},
}
