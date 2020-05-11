package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var kubernetesRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a kubernetes cluster",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("kubernetes cluster") == nil {
			kubernetesCluster, err := client.FindKubernetesCluster(args[0])
			if err != nil {
				utility.Error("Unable to find the kubernetes cluster for your search %s", err)
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
				fmt.Printf("The kubernetes cluster called %s with ID %s was delete\n", utility.Green(kubernetesCluster.Name), utility.Green(kubernetesCluster.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
