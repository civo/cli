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

var objectStoreList []utility.ObjecteList
var objectStoreDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove", "destroy"},
	Short:   "Remove an Object Store",
	Example: "civo objectstore delete OBJECTSTORE_NAME",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			objectStore, err := client.FindObjectStore(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s Object Store in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one Object Store with that name in your account")
					os.Exit(1)
				}
			}
			objectStoreList = append(objectStoreList, utility.ObjecteList{ID: objectStore.ID, Name: objectStore.Name})
		} else {
			for _, v := range args {
				objectStore, err := client.FindObjectStore(v)
				if err == nil {
					objectStoreList = append(objectStoreList, utility.ObjecteList{ID: objectStore.ID, Name: objectStore.Name})
				}
			}
		}

		objectStoreNameList := []string{}
		for _, v := range objectStoreList {
			objectStoreNameList = append(objectStoreNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(objectStoreList), "objectStore"), common.DefaultYes, strings.Join(objectStoreNameList, ", ")) {

			for _, v := range objectStoreList {
				objectStore, err := client.FindObjectStore(v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteObjectStore(objectStore.ID)
				if err != nil {
					utility.Error("Error deleting the Object Store: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range objectStoreList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("objectStore", v.Name, "Object Store")
			}

			switch common.OutputFormat {
			case "json":
				if len(objectStoreList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(objectStoreList), "objectStore"), utility.Green(strings.Join(objectStoreNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
