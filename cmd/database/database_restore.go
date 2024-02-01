package database

import (
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var (
	backup      string
	restoreName string
)

var dbRestoreCmd = &cobra.Command{
	Use:     "restore",
	Aliases: []string{"reset", "restores"},
	Short:   "Restore a database",
	Example: "civo db restore <DATABASE-NAME/ID> --name <RESTORE-NAME> --backup <BACKUP-NAME>",
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

		db, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedRestore(db.Name, common.DefaultYes, backup) {
			config := &civogo.RestoreDatabaseRequest{
				Software: db.Software,
				Name:     restoreName,
				Backup:   backup,
				Region:   client.Region,
			}
			_, err = client.RestoreDatabase(db.ID, config)
			if err != nil {
				utility.Error("%s", err)
				os.Exit(1)
			}
			ow := utility.NewOutputWriter()
			switch common.OutputFormat {
			case "json":
				ow.WriteSingleObjectJSON(common.PrettySet)
			case "custom":
				ow.WriteCustomOutput(common.OutputFields)
			default:
				fmt.Printf("Restoring database %s from from backup %s\n", utility.Green(db.Name), utility.Green(backup))
			}
		} else {
			fmt.Println("Aborted")
		}

	},
}
