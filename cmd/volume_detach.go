package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
	"time"
)

var waitVolumeDetach bool

var volumeDetachCmd = &cobra.Command{
	Use:     "detach",
	Aliases: []string{"disconnect", "unlink"},
	Example: "civo volume detach VOLUME_NAME",
	Short:   "Detach a volume",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		volume, err := client.FindVolume(args[0])
		if err != nil {
			utility.Error("Unable to find the volume for your search %s", err)
			os.Exit(1)
		}

		_, err = client.DetachVolume(volume.ID)

		if waitVolumeDetach == true {

			stillDetaching := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Detaching the volume... "
			s.Start()

			for stillDetaching {
				volumeCheck, _ := client.FindVolume(volume.ID)
				if volumeCheck.MountPoint == "" {
					stillDetaching = false
					s.Stop()
				}
				time.Sleep(5 * time.Second)
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": volume.ID, "Name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The volume called %s with ID %s was detached\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
