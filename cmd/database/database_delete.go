package database

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

var databaseList []utility.ObjecteList
var dbDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove", "destroy"},
	Short:   "Delete a database",
	Example: "civo db delete <DATABASE-NAME>",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if len(args) == 1 {
			db, err := client.FindDatabase(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s Database in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one database with that name in your account")
					os.Exit(1)
				}
			}
			databaseList = append(databaseList, utility.ObjecteList{ID: db.ID, Name: db.Name})
		} else {
			for _, v := range args {
				db, err := client.FindDatabase(v)
				if err == nil {
					databaseList = append(databaseList, utility.ObjecteList{ID: db.ID, Name: db.Name})
				}
			}
		}

		dbNameList := []string{}
		for _, v := range databaseList {
			dbNameList = append(dbNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(databaseList), "Database"), common.DefaultYes, strings.Join(dbNameList, ", ")) {

			for _, v := range databaseList {
				db, err := client.FindDatabase(v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteDatabase(db.ID)
				if err != nil {
					utility.Error("Error deleting the Database: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range databaseList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("database", v.Name, "Database")
			}

			switch common.OutputFormat {
			case "json":
				if len(databaseList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(databaseList), "database"), utility.Green(strings.Join(dbNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
