package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var volumeRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo volume rm VOLUME_NAME",
	Short:   "Remove a volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("volume") == nil {
			volume, err := client.FindVolume(args[0])
			if err != nil {
				utility.Error("Unable to find the volume for your search %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteVolume(volume.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The volume called %s with ID %s was delete\n", utility.Green(volume.Name), utility.Green(volume.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
