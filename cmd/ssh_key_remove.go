package cmd

import (
	"fmt"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var sshKeyRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Example: "civo ssh rm NAME",
	Short:   "Remove an SSH key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedDeletion("ssh key", defaultRemove) == true {
			sshKey, err := client.FindSSHKey(args[0])
			if err != nil {
				utility.Error("Finding the SSH key for your search failed with %s", err)
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
				fmt.Printf("The SSH key called %s with ID %s was deleted\n", utility.Green(sshKey.Name), utility.Green(sshKey.ID))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
