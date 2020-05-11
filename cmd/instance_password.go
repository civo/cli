package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instancePasswordCmd = &cobra.Command{
	Use:     "password",
	Short:   "Show instance's password",
	Aliases: []string{"pw"},
	Long: `Show the specified instance's default user password by part of the instance's ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* Password
	* User

Example: civo instance public-ip ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) has the password %s (and user %s)\n", utility.Green(instance.Hostname), instance.ID, utility.Green(instance.InitialPassword), utility.Green(instance.InitialUser))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Password", instance.InitialPassword)
			ow.AppendData("User", instance.InitialUser)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
