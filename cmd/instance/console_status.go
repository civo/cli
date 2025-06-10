package instance

import (
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceConsoleStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Get the status of the console for an instance",
	Example: "civo instance console status my-instance",
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

		vnc, err := client.GetInstanceVncStatus(instance.ID)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				utility.Info("The console session for instance %s (%s) does not exist or has expired.", instance.Hostname, instance.ID)
				os.Exit(0)
			}
			utility.Error("Getting console status for instance %s: %s", instance.ID, err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("instance_id", instance.ID, "Instance ID")
		ow.AppendDataWithLabel("instance_hostname", instance.Hostname, "Instance Hostname")
		ow.AppendDataWithLabel("status", vnc.Result, "Status")
		ow.AppendDataWithLabel("uri", vnc.URI, "URI")
		ow.AppendDataWithLabel("expiry", vnc.Expiration, "Expiry")

		if common.OutputFormat == "json" {
			ow.WriteSingleObjectJSON(common.PrettySet)
		} else {
			ow.WriteKeyValues()
		}
	},
}

func init() {
	instanceConsoleCmd.AddCommand(instanceConsoleStatusCmd)
}
