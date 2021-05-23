package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var sshList []utility.ObjecteList
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

		if len(args) == 1 {
			sshKey, err := client.FindSSHKey(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s SSH key in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one SSH key with that value in your account")
					os.Exit(1)
				}
			}
			sshList = append(sshList, utility.ObjecteList{ID: sshKey.ID, Name: sshKey.Name})
		} else {
			for _, v := range args {
				sshKey, err := client.FindSSHKey(v)
				if err == nil {
					sshList = append(sshList, utility.ObjecteList{ID: sshKey.ID, Name: sshKey.Name})
				}
			}
		}

		sshKeyNameList := []string{}
		for _, v := range sshList {
			sshKeyNameList = append(sshKeyNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("ssh %s", pluralize.Pluralize(len(sshList), "key")), defaultYes, strings.Join(sshKeyNameList, ", ")) {

			for _, v := range sshList {
				_, err = client.DeleteSSHKey(v.ID)
				if err != nil {
					utility.Error("error deleting the ssh key: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range sshList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch outputFormat {
			case "json":
				if len(sshList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The ssh %s (%s) has been deleted\n", pluralize.Pluralize(len(sshList), "key"), utility.Green(strings.Join(sshKeyNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
