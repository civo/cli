package cmd

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var sshKeyName string

var appSSHKeyAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"update", "new"},
	Args:    cobra.MinimumNArgs(1),
	Short:   "Update SSH Key ID's for your application",
	Long:    "\nYou can check already existing keys with `civo sshkey ls`. If no key is found, you can create one. See `civo sshkey create --help`",
	Example: "civo app sshkey add APP_NAME SSH_KEY_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		findApp, err := client.FindApplication(args[0])
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		sshKeys, err := client.ListSSHKeys()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		var found bool
		if sshKeyName != "" {
			for _, sshKey := range sshKeys {
				if sshKey.Name != sshKeyName {
					continue
				} else {
					found = true
					updateReq := &civogo.UpdateApplicationRequest{
						SSHKeyIDs: append(findApp.SSHKeyIDs, sshKey.ID),
					}
					_, err := client.UpdateApplication(findApp.ID, updateReq)
					if err != nil {
						utility.Error("%s", err)
						os.Exit(1)
					}
					fmt.Printf("Added SSH Key with name %s to application %s\n", sshKeyName, findApp.Name)
				}
			}
		} else {
			for _, sshKey := range sshKeys {
				if sshKey.ID != args[1] {
					continue
				} else {
					found = true
					updateReq := &civogo.UpdateApplicationRequest{
						SSHKeyIDs: append(findApp.SSHKeyIDs, args[1]),
					}
					_, err := client.UpdateApplication(findApp.ID, updateReq)
					if err != nil {
						utility.Error("%s", err)
						os.Exit(1)
					}
					fmt.Printf("Added SSH Key ID %s to app %s\n", args[1], findApp.Name)
				}
			}
		}
		if !found {
			utility.Error("SSH Key %s not found", args[1])
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"id": findApp.ID, "name": findApp.Name})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON(prettySet)
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			fmt.Printf("Application %s's SSH Key IDs have been updated.\n", utility.Green(findApp.Name))
		}

	},
}
