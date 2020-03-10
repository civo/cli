package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var instanceSoftRebootCmd = &cobra.Command{
	Use:   "soft-reboot",
	Short: "Spft reboot an instance",
	Long: `Nicely ask the specified instance by part of the ID or name to restart.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname

Example: civo instance soft-reboot ID/NAME`,
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

		_, err = client.SoftRebootInstance(instance.ID)
		if err != nil {
			fmt.Printf("Soft-rebooting instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being soft-rebooted\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
