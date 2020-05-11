package cmd

import (
	"fmt"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var KubernetesClusterApp string

var kubernetesAppAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"install"},
	Short:   "Add the marketplace application to a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(KubernetesClusterApp)
		if err != nil {
			fmt.Printf("Unable to find a kubernetes cluster: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Applications: args[0],
		}

		kubeCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			fmt.Printf("Unable to install the application in the kubernetes cluster: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubeCluster.ID, "Name": kubeCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The application was install in the kubernetes cluster %s\n", aurora.Green(kubeCluster.Name))
		}
	},
}
