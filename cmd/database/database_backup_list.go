package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbBackupListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Example: `civo database backup ls <DATABASE-NAME/ID>`,
	Short:   "List all database backups",
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

		backups, err := client.ListDatabaseBackup(db.ID)
		if err != nil {
			utility.Error("%s", err)
			os.Exit(1)
		}

		if backups.DatabaseID == "" {
			return
		}

		odb := utility.NewOutputWriter()
		ombk := utility.NewOutputWriter()
		osbk := utility.NewOutputWriter()
		mbk := ""
		sbk := ""
		isConfiguredScheduled := false
		isConfiguredManual := false

		odb.StartLine()
		odb.AppendDataWithLabel("database_id", utility.TrimID(backups.DatabaseID), "Database ID")
		odb.AppendDataWithLabel("database_name", backups.DatabaseName, "Database Name")
		odb.AppendDataWithLabel("software", backups.Software, "Software")

		if backups.Scheduled != nil {
			isConfiguredScheduled = true
			osbk.AppendDataWithLabel("name", backups.Scheduled.Name, "Backup Name")
			osbk.AppendDataWithLabel("schedule", backups.Scheduled.Schedule, "Schedule")
			osbk.AppendDataWithLabel("count", fmt.Sprintf("%d", backups.Scheduled.Count), "Count")

			sbk = strings.TrimSuffix(strings.Join(backups.Scheduled.Backups, ","), ",")
			osbk.AppendDataWithLabel("backups", sbk, "Backups")
		}

		if backups.Manual != nil {
			isConfiguredManual = true
			for i, m := range backups.Manual {
				if i < len(backups.Manual)-1 {
					mbk += m.Backup + ", "
				} else {
					mbk += m.Backup
				}
			}
			ombk.AppendDataWithLabel("backups", mbk, "Backups")
		}

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			odb.AppendDataWithLabel("database_id", utility.TrimID(backups.DatabaseID), "Database ID")
			odb.AppendDataWithLabel("database_name", backups.DatabaseName, "Database Name")
			odb.AppendDataWithLabel("software", backups.Software, "Software")
			odb.AppendDataWithLabel("schedule_backup_name", backups.Scheduled.Name, "Schedule Backup Name")
			odb.AppendDataWithLabel("schedule", backups.Scheduled.Schedule, "Schedule")
			odb.AppendDataWithLabel("count", fmt.Sprintf("%d", backups.Scheduled.Count), "Count")
			odb.AppendDataWithLabel("scheduled_backups", sbk, "Schedule Backups")
			odb.AppendDataWithLabel("manual_backups", mbk, "Manual Backups")
		}

		switch common.OutputFormat {
		case "json":
			odb.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			odb.WriteCustomOutput(common.OutputFields)
		default:
			if isConfiguredScheduled {
				fmt.Println("Scheduled Backups:")
			}
			osbk.WriteTable()
			if isConfiguredManual {
				fmt.Println("Manual Backups:")
			}
			ombk.WriteTable()
		}
	},
}
