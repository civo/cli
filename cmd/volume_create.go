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

var Bootable bool
var createSizeGB int

var volumeCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		volumeConfig := &civogo.VolumeConfig{
			Name:          args[0],
			SizeGigabytes: createSizeGB,
			Bootable:      Bootable,
		}

		volume, err := client.NewVolume(volumeConfig)
		if err != nil {
			fmt.Printf("Unable to create the volume: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a volume called %s with ID %s\n", aurora.Green(volume.Name), aurora.Green(volume.ID))
		}
	},
}
