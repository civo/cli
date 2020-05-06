package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var sshKeyCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new ssh key",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		// reading the file
		data, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Unable to read the ssh key file: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.NewSSHKey(args[0], string(data))
		if err != nil {
			fmt.Printf("Unable to create the ssh key: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		sshKey, err := client.FindSSHKey(args[0])
		if err != nil {
			fmt.Printf("Unable to find the ssh key: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": sshKey.ID, "Name": sshKey.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Created a ssh key called %s with ID %s\n", aurora.Green(sshKey.Name), aurora.Green(sshKey.ID))
		}
	},
}
