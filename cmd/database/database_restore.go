package database

import (
	"fmt"
	"time"

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
	Example: `civo db restore <DATABASE-NAME/ID> --name <RESTORE-NAME> --backup <BACKUP-NAME>
civo db restore <DATABASE-NAME/ID> --name <RESTORE-NAME> --target-time 2026-10-21T10:00:00Z`,
	Args: cobra.MinimumNArgs(1),
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

		if targetTime != "" {
			restoreToPointInTime(client, db)
			return
		}

		bk, err := client.GetDatabaseBackup(args[0], backup)
		if err != nil {
			utility.Error("Database backup %s", err)
			os.Exit(1)
		}

		if utility.UserConfirmedRestore(db.Name, common.DefaultYes, backup) {
			config := &civogo.RestoreDatabaseRequest{
				Name:   restoreName,
				Backup: bk.Name,
				Region: client.Region,
			}
			_, err = client.RestoreDatabase(db.ID, config)
			if err != nil {
				utility.Error("Failed to restore %s", err)
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

// restoreToPointInTime restores the database in place to --target-time. Data
// written after that moment is lost, so it asks for confirmation like the
// backup restore does.
func restoreToPointInTime(client *civogo.Client, db *civogo.Database) {
	t, err := parseTargetTime(targetTime)
	if err != nil {
		utility.Error("%s", err)
		os.Exit(1)
	}
	point := t.Format(time.RFC3339)

	if !common.DefaultYes {
		if err := utility.AskForConfirm(fmt.Sprintf("restore db %s in place to its state at %s (data written after that moment is lost)", utility.Green(db.Name), utility.Green(point))); err != nil {
			fmt.Println("Aborted")
			return
		}
	}

	_, err = client.RestoreDatabaseToPointInTime(db.ID, &civogo.RestoreDatabaseToPointInTimeRequest{
		Name:       restoreName,
		TargetTime: t,
		Region:     client.Region,
	})
	if err != nil {
		utility.Error("Failed to restore %s", err)
		os.Exit(1)
	}

	ow := utility.NewOutputWriterWithMap(map[string]string{"id": db.ID, "name": db.Name, "target_time": point})
	switch common.OutputFormat {
	case "json":
		ow.WriteSingleObjectJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		fmt.Printf("Restoring database %s to its state at %s. Data written after that moment is lost.\n", utility.Green(db.Name), utility.Green(point))
	}
}
