package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceConsoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"terminal", "shell"},
	Short:   "Get console URL for instance",
	Long: `Get the web console's URL for a given instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* URL
	* Hostname

Example: civo instance console ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s", err)
			os.Exit(1)
		}

		url, err := client.GetInstanceConsoleURL(instance.ID)
		if err != nil {
			utility.Error("Getting console URL %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has a console at %s\n", utility.Green(instance.Hostname), instance.ID,
				utility.Green(url))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("URL", url, "Console URL")
			ow.AppendData("Hostname", instance.Hostname)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
