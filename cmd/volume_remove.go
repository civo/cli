package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var volumeRemoveCmdExamples = []string{
	"civo volume rm VOLUME_NAME",
	"civo volume rm VOLUME_ID",
}

var volumeRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: strings.Join(volumeRemoveCmdExamples, "\n"),
	Short:   "Remove a volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		volume, err := client.FindVolume(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s volume in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one volume with that value in your account")
				os.Exit(1)
			}
		}

		if !utility.CanManageVolume(volume) {
			cluster, err := client.FindKubernetesCluster(volume.ClusterID)
			if err != nil {
				utility.Error("Unable to find cluster - %s", err)
				os.Exit(1)
			}

			utility.Error("Unable to %s this volume because it's being managed by your %q Kubernetes cluster", cmd.Name(), cluster.Name)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("volume", defaultYes, volume.Name) {

			_, err = client.DeleteVolume(volume.ID)
			if err != nil {
				utility.Error("error deleting the volume: %s", err)
				os.Exit(1)
			}

			ow := utility.NewOutputWriterWithMap(map[string]string{"id": volume.ID, "name": volume.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON(prettySet)
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The volume called %s with ID %s was deleted\n", utility.Green(volume.Name), utility.Green(volume.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
