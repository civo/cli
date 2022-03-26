package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var waitVolumeDetach bool

var volumeDetachCmdExamples = []string{
	"civo volume detach VOLUME_NAME",
	"civo volume detach VOLUME_ID",
}

var volumeDetachCmd = &cobra.Command{
	Use:     "detach",
	Aliases: []string{"disconnect", "unlink"},
	Example: strings.Join(volumeDetachCmdExamples, "\n"),
	Short:   "Detach a volume from an instance",
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
			utility.Error("Volume %s", err)
			os.Exit(1)
		}

		if !utility.CanManageVolume(volume) && !forceVolumeAction {
			cluster, err := client.FindKubernetesCluster(volume.ClusterID)
			if err != nil {
				utility.Error("Unable to find cluster - %s", err)
				os.Exit(1)
			}

			utility.Error("Unable to %s this volume because it's being managed by your %q Kubernetes cluster", cmd.Name(), cluster.Name)
			os.Exit(1)
		}

		_, err = client.DetachVolume(volume.ID)
		if err != nil {
			utility.Error("error detaching the volume: %s", err)
			os.Exit(1)
		}

		if waitVolumeDetach {

			stillDetaching := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Detaching the volume... "
			s.Start()

			for stillDetaching {
				volumeCheck, err := client.FindVolume(volume.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				if volumeCheck.Status == "available" {
					stillDetaching = false
					s.Stop()
				} else {
					time.Sleep(2 * time.Second)
				}
			}
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": volume.ID, "name": volume.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("The volume called %s with ID %s was detached\n", utility.Green(volume.Name), utility.Green(volume.ID))
		}
	},
}
