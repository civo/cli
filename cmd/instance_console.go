package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
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
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		url, err := client.GetInstanceConsoleURL(instance.ID)
		if err != nil {
			fmt.Printf("Getting console URL: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has a console at %s\n", aurora.Green(instance.Hostname), instance.ID,
				aurora.Green(url))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendDataWithLabel("URL", url, "Console URL")
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}
