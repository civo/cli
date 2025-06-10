package instance

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceVncStopCmd = &cobra.Command{
	Use:     "vnc-stop <INSTANCE_ID_OR_NAME>",
	Aliases: []string{"vncstop"},
	Short:   "Stop the VNC console session for an instance",
	Example: "civo instance vnc-stop my-instance",
	Args:    cobra.ExactArgs(1),
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

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s: %s", args[0], err)
			os.Exit(1)
		}

		resp, err := client.DeleteInstanceVncSession(instance.ID)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				utility.Info("There is no active VNC session for instance %s (%s) to stop.", instance.Hostname, instance.ID)
				os.Exit(0)
			}
			utility.Error("Stopping VNC session for instance %s: %s", instance.ID, err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("instance_id", instance.ID, "Instance ID")
		ow.AppendDataWithLabel("instance_hostname", instance.Hostname, "Instance Hostname")
		ow.AppendDataWithLabel("result", string(resp.Result), "Result")

		if common.OutputFormat == "human" {
			if string(resp.Result) == "ok" {
				fmt.Printf("VNC session for instance %s (%s) was stopped successfully.\n",
					utility.Green(instance.Hostname), instance.ID)
			} else {
				fmt.Printf("Failed to stop VNC session for instance %s (%s). Result: %s\n",
					utility.Red(instance.Hostname), instance.ID, string(resp.Result))
			}
		} else {
			ow.WriteSingleObjectJSON(common.PrettySet)
		}
	},
}
