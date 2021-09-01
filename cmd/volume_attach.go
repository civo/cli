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

var waitVolumeAttach bool

var volumeAttachCmdExamples = []string{
	"civo volume attach VOLUME_NAME INSTANCE_HOSTNAME",
	"civo volume attach VOLUME_ID INSTANCE_ID",
}

var volumeAttachCmd = &cobra.Command{
	Use:     "attach",
	Aliases: []string{"connect", "link"},
	Example: strings.Join(volumeAttachCmdExamples, "\n"),
	Short:   "Attach a volume to an instance",
	Args:    cobra.MinimumNArgs(2),
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

		if !utility.CanManageVolume(volume) {
			cluster, err := client.FindKubernetesCluster(volume.ClusterID)
			if err != nil {
				utility.Error("Unable to find cluster - %s", err)
				os.Exit(1)
			}

			utility.Error("Unable to %s this volume because it's being managed by your %q Kubernetes cluster", cmd.Name(), cluster.Name)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[1])
		if err != nil {
			utility.Error("Instance %s", err)
			os.Exit(1)
		}

		_, err = client.AttachVolume(volume.ID, instance.ID)
		if err != nil {
			utility.Error("error attaching the volume: %s", err)
			os.Exit(1)
		}

		if waitVolumeAttach {

			stillAttaching := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Prefix = "Attaching volume to the instance... "
			s.Start()

			for stillAttaching {
				volumeCheck, err := client.FindVolume(volume.ID)
				if err != nil {
					utility.Error("Finding the volume failed with %s", err)
					os.Exit(1)
				}
				if volumeCheck.Status == "attached" {
					stillAttaching = false
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
			fmt.Printf("The volume called %s with ID %s was attached to the instance %s\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(instance.Hostname))
		}
	},
}
