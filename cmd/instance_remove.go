package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var instanceRemoveCmd = &cobra.Command{
	Use:     "remove",
	Example: "civo instance remove ID/HOSTNAME",
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"delete", "destroy", "rm"},
	Short:   "Remove/delete instance",
	Long: `Remove the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry there is no %s instance in your account", utility.Red(args[0]))
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one instance with that name in your account")
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("instance", defaultYes, instance.Hostname) == true {

			_, err = client.DeleteInstance(instance.ID)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}

			if outputFormat == "human" {
				fmt.Printf("The instance %s (%s) has been removed\n", utility.Green(instance.Hostname), instance.ID)
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
		} else {
			fmt.Println("Operation aborted.")
		}

	},
}
