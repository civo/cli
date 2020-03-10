package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
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
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if OutputFormat == "human" {
			fmt.Printf("The instance %s (%s) has the password %s (and user %s)\n", aurora.Green(instance.Hostname), instance.ID, aurora.Green(instance.InitialPassword), aurora.Green(instance.InitialUser))
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendData("Password", instance.InitialPassword)
			ow.AppendData("User", instance.InitialUser)
			if OutputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(OutputFields)
			}
		}
	},
}
