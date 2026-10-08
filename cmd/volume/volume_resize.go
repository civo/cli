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
as soon as the new size is admitted. With --wait it follows the resize to its outcome: it exits
non-zero when the platform settles the resize without delivering the requested size, and reports
the result as unconfirmed, without failing, when the API does not report the delivered size.`,
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

		result := map[string]string{"id": volume.ID, "name": volume.Name, "size_gb": strconv.Itoa(newSizeGB), "delivered_size_gb": "", "resize_state": ""}

		if !waitVolumeResize {
			ow := utility.NewOutputWriterWithMap(result)
			switch common.OutputFormat {
			case "json":
				ow.WriteSingleObjectJSON(common.PrettySet)
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The resize of the volume called %s with ID %s to %s GB was accepted\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(strconv.Itoa(newSizeGB)))
				fmt.Println("Check `civo volume ls` to see whether the new size has been delivered")
			}
			return
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Writer = os.Stderr
		s.Prefix = fmt.Sprintf("Resizing volume %s to %d GB... ", volume.Name, newSizeGB)
		s.Start()
		outcome, err := resizeWait{
			find:      func() (*civogo.Volume, error) { return client.GetVolume(volume.ID) },
			before:    volume,
			requested: newSizeGB,
			grace:     time.Minute, timeout: 60 * time.Minute, interval: 2 * time.Second, maxFindFailures: 5,
			sleep: time.Sleep, now: time.Now,
		}.run()
		s.Stop()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}
		if v := outcome.volume; v != nil {
			if v.DeliveredSizeGigabytes > 0 {
				result["delivered_size_gb"] = strconv.Itoa(v.DeliveredSizeGigabytes)
			}
			if v.Resize != nil {
				result["resize_state"] = v.Resize.State
				if v.Resize.Reason != "" {
					result["resize_reason"] = v.Resize.Reason
				}
			}
		}
		if outcome.state == resizeFailed {
			utility.Error("The volume %s was not resized to %d GB: %s", volume.Name, newSizeGB, outcome.detail)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(result)
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			if outcome.state == resizeUnconfirmed {
				fmt.Printf("The resize of the volume called %s with ID %s to %s GB was accepted, but could not be confirmed: %s\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(strconv.Itoa(newSizeGB)), outcome.detail)
				fmt.Println("Check `civo volume ls` to see whether the new size has been delivered")
				return
			}
			fmt.Printf("The volume called %s with ID %s was resized to %s GB\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(strconv.Itoa(newSizeGB)))
		}
	},
}
