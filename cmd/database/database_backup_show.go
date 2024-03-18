package database

import (
	"os"

	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbBackupShowCmd = &cobra.Command{
	Use:     "show",
	Example: `civo database backup show <DATABASE-NAME/ID> <BACKUP-NAME/ID>`,
	Aliases: []string{"get", "inspect"},
	Args:    cobra.MinimumNArgs(2),
	Short:   "Show details of a database backup",
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		client, err := config.CivoAPIClient()
		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		bk, err := client.GetDatabaseBackup(args[0], args[1])
		if err != nil {
			utility.Error("Database backup %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriter()

		ow.StartLine()

		if !bk.IsScheduled {
			ow.AppendDataWithLabel("id", utility.TrimID(bk.ID), "ID")
		}
		ow.AppendDataWithLabel("name", bk.Name, "Name")
		ow.AppendDataWithLabel("schedule", bk.Schedule, "Schedule")

		ow.AppendDataWithLabel("database_id", utility.TrimID(bk.DatabaseID), "Database ID")
		ow.AppendDataWithLabel("database_name", bk.DatabaseName, "Database Name")

		ow.AppendDataWithLabel("software", bk.Software, "Software")
		ow.AppendDataWithLabel("status", bk.Status, "Status")

		if common.OutputFormat == "json" || common.OutputFormat == "custom" {
			ow.AppendDataWithLabel("id", bk.ID, "ID")
			ow.AppendDataWithLabel("database_id", bk.DatabaseID, "Database ID")
		}

		switch common.OutputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			ow.WriteKeyValues()
		}
	},
}
