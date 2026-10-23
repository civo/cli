package database

import (
	"fmt"
	"os"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbBackupUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"modify", "change"},
	Short:   "Update a scheduled database backup",
	Example: "civo database backup update <DATABASE-NAME/ID> --name <NEW-BACKUP-NAME> --schedule <SCHEDULE>",
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

		destination, err := backupDestinationFromFlags()
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if schedule == "" && name == "" && destination == nil {
			utility.Error("Schedule, name or a destination must be specified")
			os.Exit(1)
		}

		db, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		backupUpdateConfig := civogo.DatabaseBackupUpdateRequest{
			Region: client.Region,
		}

		if schedule != "" {
			backupUpdateConfig.Schedule = schedule
		}
		if name != "" {
			backupUpdateConfig.Name = name
		}

		var bk *civogo.DatabaseBackup
		if destination != nil {
			// Destinations exist on the v3alpha1 API only. Only the flags that
			// were given are sent: an explicit empty name or schedule would
			// clear them on the server.
			v3 := &civogo.UpdateDatabaseBackupV3alpha1Request{Region: client.Region, Destination: destination}
			if name != "" {
				v3.Name = &name
			}
			if schedule != "" {
				v3.Schedule = &schedule
			}
			bk, err = client.UpdateDatabaseBackupV3alpha1(db.ID, v3)
		} else {
			bk, err = client.UpdateDatabaseBackup(db.ID, &backupUpdateConfig)
		}
		if err != nil {
			utility.Error("Error creating database %s", err)
			os.Exit(1)
		}

		if bk.DatabaseID == "" {
			return
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"database_id":   bk.DatabaseID,
			"database_name": bk.DatabaseName,
			"software":      bk.Software,
			"backup_name":   bk.Name,
			"schedule":      bk.Schedule,
		})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Database backup (%s) for database %s has been update\n", utility.Green(bk.Name), utility.Green(bk.DatabaseName))
		}
	},
}
