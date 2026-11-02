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
	"golang.org/x/term"
)

var newSizeGB int
var waitVolumeResize bool
var waitVolumeResizeTimeout time.Duration

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
non-zero when the platform settles the resize without delivering the requested size. Exit 0 covers
two outcomes, reported in the "outcome" field: "delivered", and "unconfirmed" when no evidence of
the resize arrived within the one minute grace window or the API does not report the delivered
size. The wait gives up after --wait-timeout (60m by default) with a non-zero exit. With -o json or
-o custom the outcome, its detail and the resize state are written for every outcome, including a
failed one.`,
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

		// The API admits an equal size and writes nothing, so the platform never starts a resize
		// and --wait would sit through the grace window for nothing: refuse it up front.
		if newSizeGB <= volume.SizeGigabytes {
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

		stopSpinner := startResizeSpinner(fmt.Sprintf("Resizing volume %s to %d GB... ", volume.Name, newSizeGB))
		wait := resizeWait{
			find:      func() (*civogo.Volume, error) { return client.GetVolume(volume.ID) },
			before:    volume,
			requested: newSizeGB,
			grace:     resizeWaitGrace, timeout: waitVolumeResizeTimeout, interval: resizeWaitInterval, maxFindFailures: resizeWaitMaxFindFailures,
			sleep: time.Sleep, now: time.Now,
		}
		outcome, err := wait.run()
		stopSpinner()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}
		result["outcome"] = outcome.state
		if outcome.detail != "" {
			result["outcome_detail"] = outcome.detail
		}
		if v := outcome.volume; v != nil {
			if v.DeliveredSizeGigabytes > 0 {
				result["delivered_size_gb"] = strconv.Itoa(v.DeliveredSizeGigabytes)
			}
			// Only an outcome of this request is reported; one left over from an earlier
			// resize would be read as this one's.
			if wait.resizeFromThisRequest(v) {
				result["resize_state"] = v.Resize.State
				if v.Resize.Reason != "" {
					result["resize_reason"] = v.Resize.Reason
				}
			}
		}
		// Structured output is written for every outcome, a failed one included, so a script
		// reads why before the non-zero exit.
		ow := utility.NewOutputWriterWithMap(result)
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			switch outcome.state {
			case resizeFailed:
				utility.Error("The volume %s was not resized to %d GB: %s", volume.Name, newSizeGB, outcome.detail)
			case resizeUnconfirmed:
				fmt.Printf("The resize of the volume called %s with ID %s to %s GB was accepted, but could not be confirmed: %s\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(strconv.Itoa(newSizeGB)), outcome.detail)
				fmt.Println("Check `civo volume ls` to see whether the new size has been delivered")
			default:
				fmt.Printf("The volume called %s with ID %s was resized to %s GB\n", utility.Green(volume.Name), utility.Green(volume.ID), utility.Green(strconv.Itoa(newSizeGB)))
			}
		}
		if outcome.state == resizeFailed {
			os.Exit(1)
		}
	},
}

// startResizeSpinner shows a spinner on stderr while the wait runs, only when stderr is a terminal:
// in a pipe or a CI log the carriage-return frames are noise. It returns the function that stops it.
func startResizeSpinner(prefix string) func() {
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		return func() {}
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Prefix = prefix
	s.Start()
	return s.Stop
}
