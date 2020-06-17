package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var kubernetesRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo kubernetes remove CLUSTER_NAME",
	Short:   "Remove a Kubernetes cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("kubernetes cluster", defaultYes) == true {
			kubernetesCluster, err := client.FindKubernetesCluster(args[0])
			if err != nil {
				utility.Error("Unable to find the Kubernetes cluster for your search because of %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteKubernetesCluster(kubernetesCluster.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": kubernetesCluster.ID, "Name": kubernetesCluster.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The Kubernetes cluster called %s with ID %s was deleted\n", utility.Green(kubernetesCluster.Name), utility.Green(kubernetesCluster.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
