package cmd

import (
	"fmt"
	_ "github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var KubernetesNewName string

var kubernetesRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a kubernetes cluster",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		kubernetesFindCluster, err := client.FindKubernetesCluster(args[0])
		if err != nil {
			fmt.Printf("Unable to find a kubernetes cluster: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		configKubernetes := &civogo.KubernetesClusterConfig{
			Name: KubernetesNewName,
		}

		kubernetesCluster, err := client.UpdateKubernetesCluster(kubernetesFindCluster.ID, configKubernetes)
		if err != nil {
			fmt.Printf("Unable to rename a kubernetes cluster: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The kubernetes cluster was rename to %s with ID %s\n", aurora.Green(kubernetesCluster.Name), aurora.Green(kubernetesCluster.ID))
		}
	},
}
