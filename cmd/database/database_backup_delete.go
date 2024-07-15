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

var backupList []utility.ObjecteList
var dbBackupDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "remove", "destroy"},
	Short:   "Delete a manual database backup",
	Example: "civo database backup delete <DATABASE-NAME/ID> <BACKUP-NAME/ID>",
	Args:    cobra.MinimumNArgs(2),
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

		if len(args) == 2 {
			bk, err := client.FindDatabaseBackup(args[0], args[1])
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
			backupList = append(backupList, utility.ObjecteList{ID: bk.ID, Name: bk.Name})
		} else {
			for idx, v := range args {
				if idx == 0 {
					continue
				}
				bk, err := client.FindDatabaseBackup(args[0], v)
				if err == nil {
					backupList = append(backupList, utility.ObjecteList{ID: bk.ID, Name: bk.Name})
				}
			}
		}

		dbNameList := []string{}
		for _, v := range backupList {
			dbNameList = append(dbNameList, v.Name)
		}

		if utility.UserConfirmedDeletion(pluralize.Pluralize(len(backupList), "Database Backup"), common.DefaultYes, strings.Join(dbNameList, ", ")) {

			for _, v := range backupList {
				db, err := client.FindDatabaseBackup(args[0], v.ID)
				if err != nil {
					utility.Error("%s", err)
					os.Exit(1)
				}
				_, err = client.DeleteDatabaseBackup(args[0], db.ID)
				if err != nil {
					utility.Error("Error deleting the Database backup: %s", err)
					os.Exit(1)
				}
			}

			ow := utility.NewOutputWriter()

			for _, v := range backupList {
				ow.StartLine()
				ow.AppendDataWithLabel("id", v.ID, "ID")
				ow.AppendDataWithLabel("backup", v.Name, "Backup")
			}

			switch common.OutputFormat {
			case "json":
				if len(backupList) == 1 {
					ow.WriteSingleObjectJSON(common.PrettySet)
				} else {
					ow.WriteMultipleObjectsJSON(common.PrettySet)
				}
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("The %s (%s) has been deleted\n", pluralize.Pluralize(len(backupList), "database backup"), utility.Green(strings.Join(dbNameList, ", ")))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
