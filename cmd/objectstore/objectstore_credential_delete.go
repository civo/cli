package objectstore

import (
	"errors"
	"fmt"
	"strings"

	pluralize "github.com/alejandrojnm/go-pluralize"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var objectStoreCredsList []utility.ObjecteList
var objectStoreCredentialDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove", "destroy"},
	Short:   "Delete an Object Store Credential",
	Example: "civo objectstore credential delete CREDENTIAL-NAME",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			credential, err := client.FindObjectStoreCredential(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s Object Store credential in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one Object Store credential with that name in your account")
					os.Exit(1)
				}
			}
			objectStoreCredsList = append(objectStoreCredsList, utility.ObjecteList{ID: credential.ID, Name: credential.Name})
		} else {
			for _, v := range args {
				credential, err := client.FindObjectStoreCredential(v)
				if err == nil {
					objectStoreCredsList = append(objectStoreCredsList, utility.ObjecteList{ID: credential.ID, Name: credential.Name})
				}
			}
		}

		objectStoreCredsNameList := []string{}
		for _, v := range objectStoreCredsList {
			objectStoreCredsNameList = append(objectStoreCredsNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(objectStoreCredsList), "objectStoreCredential"), common.DefaultYes, strings.Join(objectStoreCredsNameList, ", ")) {

			for _, v := range objectStoreCredsList {
				credential, err := client.FindObjectStoreCredential(v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteObjectStoreCredential(credential.ID)
				if err != nil {
					utility.Error("Error deleting the Object Store Credential: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range objectStoreCredsList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("objectStoreCredential", v.Name, "Object Store Credential")
			}

			switch common.OutputFormat {
			case "json":
				if len(objectStoreCredsList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(objectStoreCredsList), "objectStoreCredential"), utility.Green(strings.Join(objectStoreCredsNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
