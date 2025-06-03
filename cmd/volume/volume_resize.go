package volume

import (
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var newSizeGB int

var volumeResizeCmd = &cobra.Command{
	Use:     "resize",
	Short:   "Resize a volume",
	Example: "civo volume resize VOLUME_NAME --size-gb=100",
	Args:    cobra.MinimumNArgs(1),
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

		volume, err := client.FindVolume(args[0])
		if err != nil {
			utility.Error("Volume %s", err)
			os.Exit(1)
		}

		if newSizeGB < volume.SizeGigabytes {
			utility.Printf("Sorry, the volume size specified (%s) must be larger than the volume's current size (%s)\n", utility.Red(strconv.Itoa(newSizeGB)), utility.Green(strconv.Itoa(volume.SizeGigabytes)))
			os.Exit(1)
		}

		_, err = client.ResizeVolume(volume.ID, newSizeGB)
		if err != nil {
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
			utility.Printf("The volume called %s with ID %s was resized\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
