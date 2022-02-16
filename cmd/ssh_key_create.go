package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var keyCreate string

var sshKeyCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo ssh create NAME --key PATH_TO_SSH_KEY",
	Short:   "Create a new SSH key",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		// reading the file
		data, err := ioutil.ReadFile(keyCreate)
		if err != nil {
			utility.Error("Reading the SSH key file failed with %s", err)
			os.Exit(1)
		}

		//validate the ssh public key
		if err := utility.ValidateSSHKey(data); err != nil {
			utility.Error("Validating the SSH key failed with %s", err)
			os.Exit(1)
		}

		_, err = client.NewSSHKey(args[0], string(data))
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		sshKey, err := client.FindSSHKey(args[0])
		if err != nil {
			utility.Error("SSHKey %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": sshKey.ID, "name": sshKey.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created an SSH key called %s with ID %s\n", utility.Green(sshKey.Name), utility.Green(sshKey.ID))
		}
	},
}
