package cmd

import (
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var sshKeyListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo ssh ls`,
	Short:   "List all SSH keys",
	Long: `List all current SSH keys.
If you wish to use a custom format, the available fields are:

	* ID
	* Name
	* Fingerprint

Example: civo ssh ls -o custom -f "ID: Name"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		sshKeys, err := client.ListSSHKeys()
		if err != nil {
			utility.Error("%s", err)
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
