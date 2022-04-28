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

var appList []utility.ObjecteList
var appRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"delete", "rm"},
	Short:   "Remove an application",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo app rm APP_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}
		if len(args) == 1 {
			app, err := client.FindApplication(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s application in your account", utility.Red(args[0]))
					os.Exit(1)
				} else if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one application with that name in your account")
					os.Exit(1)
				} else {
					utility.Error("%s", err)
					os.Exit(1)
				}
			}
			appList = append(appList, utility.ObjecteList{ID: app.ID, Name: app.Name})
		} else {
			for _, v := range args {
				app, err := client.FindApplication(v)
				if err == nil {
					appList = append(appList, utility.ObjecteList{ID: app.ID, Name: app.Name})
				}
			}
		}

		appNameList := []string{}
		for _, v := range appList {
			appNameList = append(appNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("application %s", pluralize.Pluralize(len(appList), "")), defaultYes, strings.Join(appNameList, ", ")) {

			for _, v := range appList {
				_, err = client.DeleteApplication(v.ID)
				if err != nil {
					utility.Error("error deleting application: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range appList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch outputFormat {
			case "json":
				if len(appList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The application %s(%s) has been deleted\n", pluralize.Pluralize(len(appList), ""), utility.Green(strings.Join(appNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
