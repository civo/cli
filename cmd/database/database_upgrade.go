package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/spf13/cobra"
)

var upgradeTargetVersion string
var upgradeDryRun bool
var upgradePostgresVersion string

var dbUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade a database's engine release or its PostgreSQL major version",
	Long: `Without --postgres-version, moves a PostgreSQL database to a newer engine
release, which brings newer PostgreSQL versions and fixes. Databases on the same
network share the engine, so they are upgraded together and each restarts while
it moves.

With --postgres-version, upgrades this one database in place to a newer
PostgreSQL major version, keeping its data, endpoint and credentials. A backup
named pre-upgrade-pgN is taken first and kept, the database stops for a few
minutes while its data is upgraded, its replicas re-sync afterwards, and
point-in-time restore starts again from the first backup on the new version.
There is no downgrade.

Nothing is upgraded until you run this command. "civo db show" reports the
engine release, the engine upgrade available and the PostgreSQL majors the
database can be upgraded to.`,
	Example: `civo db upgrade <DATABASE-NAME/ID> --dry-run
civo db upgrade <DATABASE-NAME/ID>
civo db upgrade <DATABASE-NAME/ID> --version 2.9.0 --yes
civo db upgrade <DATABASE-NAME/ID> --postgres-version 18 --dry-run
civo db upgrade <DATABASE-NAME/ID> --postgres-version 18`,
	Args: cobra.ExactArgs(1),
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

		if upgradePostgresVersion != "" {
			if upgradeTargetVersion != "" {
				utility.Error("--version picks an engine release and --postgres-version a PostgreSQL major; use one at a time")
				os.Exit(1)
			}
			runMajorVersionUpgrade(client, db)
			return
		}

		// Always ask the API what would be upgraded first, so the prompt can
		// name every database that will restart.
		plan, err := client.UpgradeDatabase(db.ID, &civogo.UpgradeDatabaseRequest{TargetVersion: upgradeTargetVersion, ValidateOnly: true})
		if err != nil {
			utility.Error("Failed to plan the upgrade: %s", err)
			os.Exit(1)
		}

		if upgradeDryRun {
			writeUpgrade(db, plan)
			return
		}

		names := make([]string, 0, len(plan.Databases))
		for _, d := range plan.Databases {
			names = append(names, d.Name)
		}
		if !common.DefaultYes {
			msg := fmt.Sprintf("upgrade the engine from %s to %s for %s (each restarts during the upgrade)", utility.Green(plan.FromVersion), utility.Green(plan.ToVersion), utility.Green(strings.Join(names, ", ")))
			if err := utility.AskForConfirm(msg); err != nil {
				fmt.Println("Aborted")
				return
			}
		}

		result, err := client.UpgradeDatabase(db.ID, &civogo.UpgradeDatabaseRequest{TargetVersion: plan.ToVersion})
		if err != nil {
			utility.Error("Failed to upgrade the database: %s", err)
			os.Exit(1)
		}
		writeUpgrade(db, result)
	},
}

func writeUpgrade(db *civogo.Database, u *civogo.DatabaseUpgrade) {
	names := make([]string, 0, len(u.Databases))
	for _, d := range u.Databases {
		names = append(names, d.Name)
	}
	ow := utility.NewOutputWriterWithMap(map[string]string{
		"id":           db.ID,
		"name":         db.Name,
		"from_version": u.FromVersion,
		"to_version":   u.ToVersion,
		"databases":    strings.Join(names, ","),
		"started":      fmt.Sprintf("%t", u.Started),
	})
	switch common.OutputFormat {
	case "json":
		ow.WriteSingleObjectJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		if u.Started {
			fmt.Printf("Upgrading the engine from %s to %s for %s. Each database restarts during the upgrade; follow progress with civo db show.\n",
				utility.Green(u.FromVersion), utility.Green(u.ToVersion), utility.Green(strings.Join(names, ", ")))
		} else {
			fmt.Printf("Upgrading %s would move the engine from %s to %s and restart: %s\n",
				utility.Green(db.Name), utility.Green(u.FromVersion), utility.Green(u.ToVersion), strings.Join(names, ", "))
		}
	}
}
