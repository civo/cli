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

var dbeList []utility.ObjectList
var dbDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "destroy"},
	Example: "civo db delete DATABASE_NAME",
	Short:   "Delete a database",
	Args:    cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return getAllDatabasesList(), cobra.ShellCompDirectiveNoFileComp
		}
		return getDatabasesWithSimilarParam(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if regionSet != "" {
			client.Region = regionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if len(args) == 1 {
			db, err := client.FindDatabase(args[0])
			if err != nil {
				if errors.Is(err, civogo.ZeroMatchesError) {
					utility.Error("sorry there is no %s database in your account", utility.Red(args[0]))
					os.Exit(1)
				}
				if errors.Is(err, civogo.MultipleMatchesError) {
					utility.Error("sorry we found more than one database with that name in your account")
					os.Exit(1)
				}
			}
			dbeList = append(dbeList, utility.ObjectList{ID: db.ID, Name: db.Name})
		} else {
			for _, v := range args {
				db, err := client.FindDatabase(v)
				if err == nil {
					dbeList = append(dbeList, utility.ObjectList{ID: db.ID, Name: db.Name})
				}
			}
		}

		dbNameList := []string{}
		for _, v := range dbeList {
			dbNameList = append(dbNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(fmt.Sprintf("Database %s", pluralize.Pluralize(len(dbeList), "")), defaultYes, strings.Join(dbNameList, ", ")) {

			for _, v := range dbeList {
				_, err = client.DeleteDatabase(v.ID)
				if err != nil {
					utility.Error("error deleting the database: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range dbeList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("name", v.Name, "Name")
			}

			switch outputFormat {
			case "json":
				if len(dbeList) == 1 {
					ow.WriteSingleObjectJSON(prettySet)
				} else {
					ow.WriteMultipleObjectsJSON(prettySet)
				}
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The Database %s(%s) has been deleted\n", pluralize.Pluralize(len(dbeList), ""), utility.Green(strings.Join(dbNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted.")
		}
	},
}
