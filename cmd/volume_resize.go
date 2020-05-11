package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var newSizeGB int

var volumeResizeCmd = &cobra.Command{
	Use:   "resize",
	Short: "Resize a volume",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		volume, err := client.FindVolume(args[0])
		if err != nil {
			utility.Error("Unable to find the volume %s", err)
			os.Exit(1)
		}

		if newSizeGB < volume.SizeGigabytes {
			fmt.Printf("Sorry, The volume size specified (%s) wasn't larger than the volume's current size (%s)\n", utility.Red(strconv.Itoa(newSizeGB)), utility.Green(strconv.Itoa(volume.SizeGigabytes)))
			os.Exit(1)
		}

		_, err = client.ResizeVolume(volume.ID, newSizeGB)
		if err != nil {
			utility.Error("Unable to resize the volume %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The a volume called %s with ID %s was resized\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
