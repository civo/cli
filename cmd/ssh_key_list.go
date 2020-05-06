package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
)

var sshKeyListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List all ssh keys",
	Long: `List all current ssh keys.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Fingerprint

Example: civo ssh ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		sshKeys, err := client.ListSSHKeys()
		if err != nil {
			fmt.Printf("Unable to list ssh keys: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()
		for _, sshkey := range sshKeys {
			ow.StartLine()

			ow.AppendData("ID", sshkey.ID)
			ow.AppendData("Name", sshkey.Name)
			ow.AppendDataWithLabel("Fingerprint", sshkey.Fingerprint, "Finger Print")
		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
