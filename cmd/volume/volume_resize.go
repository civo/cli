package volume

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var newSizeGB int
var waitVolumeResize bool

var volumeResizeCmd = &cobra.Command{
	Use:     "resize",
	Short:   "Resize a volume",
	Example: "civo volume resize VOLUME_NAME --size-gb=100 --wait",
	Long: `Grow a volume to a new size. The size is additive only.

A detached volume is resized offline. An attached volume is grown in place when its volume type
supports online expansion and the instance is running; otherwise the API refuses the request and
says what to do (detach the volume, start the instance, or move the volume to another volume type).

The API accepts the request before the platform carries it out: without --wait the command returns
as soon as the new size is admitted. With --wait it follows the resize to its outcome and exits
non-zero when the volume does not end up delivering the requested size.`,
	Args: cobra.MinimumNArgs(1),
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
			fmt.Printf("Sorry, the volume size specified (%s) must be larger than the volume's current size (%s)\n", utility.Red(strconv.Itoa(newSizeGB)), utility.Green(strconv.Itoa(volume.SizeGigabytes)))
			os.Exit(1)
		}

		_, err = client.ResizeVolume(volume.ID, newSizeGB)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		result := map[string]string{"id": volume.ID, "name": volume.Name, "size_gb": strconv.Itoa(newSizeGB)}
		verb := "is being resized"

		if waitVolumeResize {
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Writer = os.Stderr
			s.Prefix = fmt.Sprintf("Resizing volume %s to %d GB... ", volume.Name, newSizeGB)
			s.Start()
			settled, err := waitForResize(func() (*civogo.Volume, error) { return client.GetVolume(volume.ID) }, newSizeGB, time.Minute, 60*time.Minute, 2*time.Second, time.Sleep, time.Now)
			s.Stop()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			if settled.DeliveredSizeGigabytes > 0 {
				result["delivered_size_gb"] = strconv.Itoa(settled.DeliveredSizeGigabytes)
			}
			if settled.Resize != nil {
				result["resize_state"] = settled.Resize.State
				result["resize_reason"] = settled.Resize.Reason
			}
			if outcome := resizeOutcome(settled, newSizeGB); outcome != "" {
				utility.Error("The volume %s was not resized to %d GB: %s", volume.Name, newSizeGB, outcome)
				os.Exit(1)
			}
			verb = "was resized"
		}

		ow := utility.NewOutputWriterWithMap(result)

		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("The volume called %s with ID %s %s to %s GB\n", utility.Green(volume.Name), utility.Green(volume.ID), verb, utility.Green(strconv.Itoa(newSizeGB)))
			if !waitVolumeResize {
				fmt.Println("Run again with --wait, or check `civo volume ls`, to see whether the new size was delivered")
			}
		}
	},
}
