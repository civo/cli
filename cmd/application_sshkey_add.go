package cmd

import (
	"fmt"
	"os/exec"
	"strings"

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

		if sshKeyName != "" {
			sshKey, err := client.FindSSHKey(sshKeyName)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			sshKeyLsCmd, err := exec.Command("civo", "sshkey", "ls").Output()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			} else {
				sshKeyList := strings.Split(string(sshKeyLsCmd), " ")
				if contains(sshKeyList, sshKey.ID) {

					updateReq := &civogo.UpdateApplicationRequest{
						SSHKeyIDs: append(findApp.SSHKeyIDs, sshKey.ID),
					}

					_, err := client.UpdateApplication(findApp.ID, updateReq)
					if err != nil {
						utility.Error("%s", err)
						os.Exit(1)
					}
					fmt.Printf("Added SSH Key with name %s to application %s\n", sshKey.Name, findApp.Name)
				} else {
					utility.Error("SSH Key with name %s not found", sshKey.Name)
					os.Exit(1)
				}
			}
		} else {
			sshKeyLsCmd, err := exec.Command("civo", "sshkey", "ls").Output()
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			} else {
				sshKeyList := strings.Split(string(sshKeyLsCmd), " ")
				if contains(sshKeyList, args[1]) {

					updateReq := &civogo.UpdateApplicationRequest{
						SSHKeyIDs: append(findApp.SSHKeyIDs, args[1]),
					}

					_, err := client.UpdateApplication(findApp.ID, updateReq)
					if err != nil {
						utility.Error("%s", err)
						os.Exit(1)
					}
					fmt.Printf("Added SSH Key ID %s to app %s\n", args[1], args[0])
				} else {
					utility.Error("SSH Key ID %s not found", args[1])
					os.Exit(1)
				}
			}
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
