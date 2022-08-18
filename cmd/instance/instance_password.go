package instance

import (
	"fmt"
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instancePasswordCmd = &cobra.Command{
	Use:     "password",
	Example: "civo instance public-ip ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Show instance's default password",
	Aliases: []string{"pw"},
	Long: `Show the specified instance's default SSH password by part of the instance's ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* Password
	* User`,
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
			utility.Error("%s", err)
			os.Exit(1)
		}

		if common.OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has the password %s (and user %s)\n", utility.Green(instance.Hostname), instance.ID, utility.Green(instance.InitialPassword), utility.Green(instance.InitialUser))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendDataWithLabel("ID", instance.ID, "")
			ow.AppendDataWithLabel("Hostname", instance.Hostname, "")
			ow.AppendDataWithLabel("Password", instance.InitialPassword, "")
			ow.AppendDataWithLabel("User", instance.InitialUser, "")
			if common.OutputFormat == "json" {
				ow.WriteSingleObjectJSON(common.PrettySet)
			} else {
				ow.WriteCustomOutput(common.OutputFields)
			}
		}
	},
}
