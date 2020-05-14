package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceRemoveCmd = &cobra.Command{
	Use:     "remove",
	Example: "civo instance remove ID/HOSTNAME",
	Aliases: []string{"delete", "destroy", "rm"},
	Short:   "Remove/delete instance",
	Long: `Remove the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("instance") == nil {
			instance, err := client.FindInstance(args[0])
			if err != nil {
				utility.Error("Finding instance %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteInstance(instance.ID)
			if err != nil {
				utility.Error("Removing instance %s", err)
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
