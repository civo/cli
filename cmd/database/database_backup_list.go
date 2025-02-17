package database

import (
	"os"

	"github.com/civo/civogo"
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

		switch db.Software {
		case "PostgreSQL":
			postgresScheduledBackups(backups)
			postgresManualBackups(backups)
		case "MySQL":
			mysqlBackups(backups)
		}
	},
}

func postgresScheduledBackups(backups *civogo.PaginatedDatabaseBackup) {
	ow := utility.NewOutputWriter()
	printMsg := false
	for _, bk := range backups.Items {
		if !bk.IsScheduled {
			continue
		}
		printMsg = true
		ow.StartLine()
		ow.AppendDataWithLabel("name", bk.Name, "Name")
		ow.AppendDataWithLabel("schedule", bk.Schedule, "Schedule")

		ow.AppendDataWithLabel("database_id", utility.TrimID(bk.DatabaseID), "Database ID")
		ow.AppendDataWithLabel("database_name", bk.DatabaseName, "Database Name")

		ow.AppendDataWithLabel("software", bk.Software, "Software")
		ow.AppendDataWithLabel("status", bk.Status, "Status")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("database_id", bk.DatabaseID, "Database ID")
		}
	}

	switch common.OutputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		if printMsg {
			utility.Println("Scheduled backup")
		}
		ow.WriteTable()
	}
}

func postgresManualBackups(backups *civogo.PaginatedDatabaseBackup) {
	ow := utility.NewOutputWriter()
	printMsg := false
	for _, bk := range backups.Items {
		if bk.IsScheduled {
			continue
		}
		printMsg = true
		ow.StartLine()
		ow.AppendDataWithLabel("id", utility.TrimID(bk.ID), "ID")
		ow.AppendDataWithLabel("name", bk.Name, "Name")

		ow.AppendDataWithLabel("database_id", utility.TrimID(bk.DatabaseID), "Database ID")
		ow.AppendDataWithLabel("database_name", bk.DatabaseName, "Database Name")

		ow.AppendDataWithLabel("software", bk.Software, "Software")
		ow.AppendDataWithLabel("status", bk.Status, "Status")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("id", bk.ID, "ID")
			ow.AppendDataWithLabel("database_id", bk.DatabaseID, "Database ID")
		}
	}

	switch common.OutputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		if printMsg {
			utility.Println("Manual backups")
		}
		ow.WriteTable()
	}
}

func mysqlBackups(backups *civogo.PaginatedDatabaseBackup) {
	ow := utility.NewOutputWriter()
	printMsg := false
	for _, bk := range backups.Items {
		printMsg = true
		ow.StartLine()
		ow.AppendDataWithLabel("id", utility.TrimID(bk.ID), "ID")
		ow.AppendDataWithLabel("name", bk.Name, "Name")

		ow.AppendDataWithLabel("database_id", utility.TrimID(bk.DatabaseID), "Database ID")
		ow.AppendDataWithLabel("database_name", bk.DatabaseName, "Database Name")

		ow.AppendDataWithLabel("software", bk.Software, "Software")
		ow.AppendDataWithLabel("status", bk.Status, "Status")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("id", bk.ID, "ID")
			ow.AppendDataWithLabel("database_id", bk.DatabaseID, "Database ID")
		}
	}

	switch common.OutputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		if printMsg {
			utility.Println("Manual backups")
		}
		ow.WriteTable()
	}
}
