package database

import (
	"fmt"
	"os"

	"github.com/adhocore/gronx"
	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var dbBackupCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Example: "civo database backup create <DATABASE-NAME/ID> --name <BACKUP_NAME> --schedule <SCHEDULE> --count <COUNT>",
	Short:   "Create a new database backup",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utility.EnsureCurrentRegion()

		check, region, err := utility.CheckAvailability("dbaas", common.RegionSet)
		if err != nil {
			utility.Error("Error checking availability %s", err)
			os.Exit(1)
		}

		if !check {
			utility.Error("Sorry you can't create a database in the %s region", region)
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		if common.RegionSet != "" {
			client.Region = common.RegionSet
		}

		if count <= 0 {
			utility.Error("Count must be greater than zero, you have given: %d", count)
			os.Exit(1)
		}

		if schedule == "" {
			utility.Error("Schedule must be specified")
			os.Exit(1)
		}

		gron := gronx.New()
		if !gron.IsValid(schedule) {
			utility.Error("Invalid schedule, you have given: %s", schedule)
			os.Exit(1)
		}

		db, err := client.FindDatabase(args[0])
		if err != nil {
			utility.Error("Database %s", err)
			os.Exit(1)
		}

		backupCreateConfig := civogo.DatabaseBackupCreateRequest{
			Name:     name,
			Schedule: schedule,
			Count:    int32(count),
			Region:   client.Region,
		}

		bk, err := client.CreateDatabaseBackup(db.ID, &backupCreateConfig)
		if err != nil {
			utility.Error("Error creating database %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{
			"database_id":   bk.DatabaseID,
			"database_name": bk.DatabaseName,
			"software":      bk.Software,
			"name":          bk.Name,
			"schedule":      bk.Schedule,
			"count":         fmt.Sprintf("%d", count),
		})
		switch common.OutputFormat {
		case "json":
			ow.WriteSingleObjectJSON(common.PrettySet)
		case "custom":
			ow.WriteCustomOutput(common.OutputFields)
		default:
			fmt.Printf("Database backup (%s) for database %s has been created\n", utility.Green(bk.Name), utility.Green(bk.DatabaseName))
		}
	},
}
