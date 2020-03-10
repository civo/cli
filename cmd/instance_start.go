package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var instanceStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"boot", "run"},
	Short:   "Start an instance",
	Long: `Pull the power and restart the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname

Example: civo instance reboot ID/NAME`,
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

		_, err = client.StartInstance(instance.ID)
		if err != nil {
			fmt.Printf("Starting instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being started\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}
