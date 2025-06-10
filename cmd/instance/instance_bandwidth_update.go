package instance

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var bandwidthLimitUpdate int

var instanceBandwidthUpdateCmd = &cobra.Command{
	Use:     "bandwidth-update <INSTANCE_ID_OR_NAME>",
	Aliases: []string{"update-bandwidth"},
	Short:   "Update the network bandwidth limit for an instance",
	Long: `Update the network bandwidth limit for a specified instance.
The limit is specified in Mbps. Use 0 for unlimited bandwidth (if supported by the API).`,
	Example: "civo instance bandwidth-update my-instance --limit 1000",
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

		resp, err := client.UpdateInstanceBandwidth(instance.ID, bandwidthLimitUpdate)
		if err != nil {
			utility.Error("Updating bandwidth limit for instance %s: %s", instance.ID, err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("id", instance.ID, "ID")
		ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
		ow.AppendDataWithLabel("result", string(resp.Result), "Result")

		newLimitStr := strconv.Itoa(bandwidthLimitUpdate) + " Mbps"
		if bandwidthLimitUpdate == 0 {
			newLimitStr = "Unlimited"
		}
		ow.AppendDataWithLabel("new_bandwidth_limit", newLimitStr, "New Bandwidth Limit")

		if common.OutputFormat == "human" {
			if string(resp.Result) == "success" {
				fmt.Printf("Network bandwidth limit for instance %s (%s) updated successfully to %s.\n",
					utility.Green(instance.Hostname), instance.ID, newLimitStr)
			} else {
				fmt.Printf("Failed to update network bandwidth limit for instance %s (%s). Result: %s\n",
					utility.Red(instance.Hostname), instance.ID, string(resp.Result))
			}
		} else {
			ow.WriteSingleObjectJSON(common.PrettySet)
		}
	},
}
