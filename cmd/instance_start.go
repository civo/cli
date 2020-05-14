package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceStartCmd = &cobra.Command{
	Use:     "start",
	Example: "civo instance reboot ID/HOSTNAME",
	Aliases: []string{"boot", "run"},
	Short:   "Start an instance",
	Long: `Pull the power and restart the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname`,
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

		_, err = client.StartInstance(instance.ID)
		if err != nil {
			utility.Error("Starting instance %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being started\n", utility.Green(instance.Hostname), instance.ID)
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
