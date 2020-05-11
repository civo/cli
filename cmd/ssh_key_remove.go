package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"os"
)

var sshKeyRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a ssh key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		if utility.AskForConfirmDelete("ssh key") == nil {
			sshKey, err := client.FindSSHKey(args[0])
			if err != nil {
				utility.Error("Unable to find ssh key for your search %s", err)
				os.Exit(1)
			}

			_, err = client.DeleteSSHKey(sshKey.ID)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": sshKey.ID, "Name": sshKey.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The ssh key called %s with ID %s was delete\n", utility.Green(sshKey.Name), utility.Green(sshKey.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
