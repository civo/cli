package volume

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

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

		check, region, err := utility.CheckAvailability("volume", common.RegionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if !check {
			utility.Error("Sorry you can't create a volume in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		volumeConfig := &civogo.VolumeConfig{
			Name:          args[0],
			SizeGigabytes: createSizeGB,
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
			if err == civogo.QuotaLimitReachedError {
				utility.Info("Please consider deleting dangling volumes, if any. To check if you have any dangling volumes, run `civo volume ls --dangling`")
				os.Exit(1)
			}
			utility.Error("%s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": volume.ID, "name": volume.Name})

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Created a volume called %s with ID %s\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
