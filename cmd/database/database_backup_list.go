package database

import (
	"fmt"
	"os"

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

		ow := utility.NewOutputWriter()
		ow.StartLine()
		ow.AppendDataWithLabel("database_id", utility.TrimID(backups.DatabaseID), "Database ID")
		ow.AppendDataWithLabel("database_name", backups.DatabaseName, "Database Name")
		ow.AppendDataWithLabel("name", backups.Name, "Backup Name")
		ow.AppendDataWithLabel("schedule", backups.Schedule, "Schedule")
		ow.AppendDataWithLabel("count", fmt.Sprintf("%d", backups.Count), "Count")

		bk := ""
		for i, backup := range backups.Backups {
			bk += backup
			if i < len(backups.Backups)-1 {
				bk += "\n"
			}
		}
		ow.AppendDataWithLabel("backups", bk, "Backups")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("software", backups.Software, "Software")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteTable()
		}
	},
}
