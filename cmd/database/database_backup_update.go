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
	Short:   "Update a database backup",
	Example: "civo database backup update <DATABASE-NAME/ID> --name <NEW_BACKUP-NAME> --schedule <SCHEDULE> --count <COUNT>",
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

		if schedule == "" && count == 0 && name == "" {
			utility.Error("Schedule, name or count must be specified")
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
		if count >= 0 {
			backupUpdateConfig.Count = int32(count)
		}
		if name != "" {
			backupUpdateConfig.Name = name
		}

		bk, err := client.UpdateDatabaseBackup(db.ID, &backupUpdateConfig)
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
			"name":          bk.Scheduled.Name,
			"schedule":      bk.Scheduled.Schedule,
			"count":         fmt.Sprintf("%d", count),
		})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Database backup (%s) for database %s has been update\n", utility.Green(bk.Scheduled.Name), utility.Green(bk.DatabaseName))
		}
	},
}
