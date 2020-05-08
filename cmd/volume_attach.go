package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var waitVolumeAttach bool

var volumeAttachCmd = &cobra.Command{
	Use:     "attach",
	Aliases: []string{"connect", "link"},
	Short:   "Attach a volume",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		volume, err := client.FindVolume(args[0])
		if err != nil {
			fmt.Printf("Unable to find the volume for your search: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[1])
		if err != nil {
			fmt.Printf("Unable to find the instance for your search: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.AttachVolume(volume.ID, instance.ID)

		if waitVolumeAttach == true {

			stillAttaching := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = " Attaching volume to the instance..."
			s.Start()

			for stillAttaching {
				volumeCheck, _ := client.FindVolume(volume.ID)
				if volumeCheck.MountPoint != "" {
					stillAttaching = false
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
			fmt.Printf("The volume called %s with ID %s was attached to the instance %s\n", aurora.Green(volume.Name), aurora.Green(volume.ID), aurora.Green(instance.Hostname))
		}
	},
}
