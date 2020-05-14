package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var keyCreate string

var sshKeyCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo ssh create NAME --key PATH_TO_SSH_KEY",
	Short:   "Create a new ssh key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s", err)
			os.Exit(1)
		}

		// reading the file
		data, err := ioutil.ReadFile(keyCreate)
		if err != nil {
			utility.Error("Unable to read the ssh key file %s", err)
			os.Exit(1)
		}

		_, err = client.NewSSHKey(args[0], string(data))
		if err != nil {
			utility.Error("Unable to create the ssh key %s", err)
			os.Exit(1)
		}

		sshKey, err := client.FindSSHKey(args[0])
		if err != nil {
			utility.Error("Unable to find the ssh key %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": sshKey.ID, "Name": sshKey.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a ssh key called %s with ID %s\n", utility.Green(sshKey.Name), utility.Green(sshKey.ID))
		}
	},
}
