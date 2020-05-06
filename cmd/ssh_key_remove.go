package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
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
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		sshKey, err := client.FindSSHKey(args[0])
		if err != nil {
			fmt.Printf("Unable to find ssh key for your search: %s\n", aurora.Red(err))
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
			fmt.Printf("The ssh key called %s with ID %s was delete\n", aurora.Green(sshKey.Name), aurora.Green(sshKey.ID))
		}
	},
}
