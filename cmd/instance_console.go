package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceConsoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"terminal", "shell"},
	Args:    cobra.MinimumNArgs(1),
	Example: "civo instance console HOSTNAME/INSTANCE_ID",
	Short:   "Get console URL for instance",
	Long: `Get the web console's URL for a given instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* id
	* url
	* hostname

Example: civo instance console ID/NAME`,
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
			utility.Error("Instance %s", err)
			os.Exit(1)
		}

		url, err := client.GetInstanceConsoleURL(instance.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has a console at %s\n", utility.Green(instance.Hostname), instance.ID,
				utility.Green(url))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("url", url, "Console URL")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
