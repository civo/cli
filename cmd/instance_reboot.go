package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceRebootCmd = &cobra.Command{
	Use:     "reboot",
	Example: "civo instance reboot ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"hard-reboot"},
	Short:   "Hard reboot an instance",
	Long: `Pull the power and restart the specified instance by part of its ID or name.
If you wish to use a custom format, the available fields are:

	* id
	* hostname`,
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

		_, err = client.RebootInstance(instance.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being rebooted\n", utility.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("id", instance.ID, "ID")
			ow.AppendDataWithLabel("hostname", instance.Hostname, "Hostname")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
