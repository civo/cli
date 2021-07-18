package cmd

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var bootableVolume bool
var createSizeGB int
var networkVolumeID string

var volumeCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo volume create NAME [flags]",
	Short:   "Create a new volume",
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

		volumeConfig := &civogo.VolumeConfig{
			Name:          args[0],
			SizeGigabytes: createSizeGB,
			Bootable:      bootableVolume,
			Region:        client.Region,
		}

		if networkVolumeID == "default" {
			network, err := client.GetDefaultNetwork()
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}

			volumeConfig.NetworkID = network.ID

		} else {
			network, err := client.FindNetwork(networkVolumeID)
			if err != nil {
				utility.Error("Network %s", err)
				os.Exit(1)
			}

			volumeConfig.NetworkID = network.ID
		}

		volume, err := client.NewVolume(volumeConfig)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": volume.ID, "name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a volume called %s with ID %s\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
