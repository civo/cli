package database

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/utility"
)

// checkUpgradeFlags rejects combinations of the three kinds of upgrade.
func checkUpgradeFlags(engineVersion, postgresVersion string, minor bool) error {
	if minor && engineVersion != "" {
		return errors.New("--version picks an engine release and --minor the newest PostgreSQL minor; use one at a time")
	}
	if minor && postgresVersion != "" {
		return errors.New("--minor upgrades within the PostgreSQL major and --postgres-version to a new major; use one at a time")
	}
	return nil
}

// runMinorVersionUpgrade upgrades one PostgreSQL database to the newest
// minor of its major. The API is always asked first, so a dry run and the
// confirmation name the exact versions, and a refusal is shown before
// anything changes.
func runMinorVersionUpgrade(client *civogo.Client, db *civogo.Database) {
	plan, err := client.UpgradeDatabaseMinorVersion(db.ID, &civogo.UpgradeDatabaseMinorVersionRequest{ValidateOnly: true})
	if err != nil {
		utility.Error("%s cannot be upgraded: %s", db.Name, err)
		os.Exit(1)
	}

	if upgradeDryRun {
		writeMinorVersionUpgrade(db, plan)
		return
	}

	if !common.DefaultYes {
		fmt.Println(minorVersionUpgradeWarning(plan.FromVersion, plan.ToVersion))
		msg := fmt.Sprintf("upgrade %s from PostgreSQL %s to %s", utility.Green(db.Name), utility.Green(plan.FromVersion), utility.Green(plan.ToVersion))
		if err := utility.AskForConfirm(msg); err != nil {
			fmt.Println("Aborted")
			return
		}
	}

	result, err := client.UpgradeDatabaseMinorVersion(db.ID, nil)
	if err != nil {
		utility.Error("Failed to upgrade the database: %s", err)
		os.Exit(1)
	}
	writeMinorVersionUpgrade(db, result)
}

func writeMinorVersionUpgrade(db *civogo.Database, u *civogo.UpgradeDatabaseMinorVersionResponse) {
	ow := utility.NewOutputWriterWithMap(map[string]string{
		"id":           db.ID,
		"name":         db.Name,
		"from_version": u.FromVersion,
		"to_version":   u.ToVersion,
		"started":      fmt.Sprintf("%t", u.Started),
	})
	switch common.OutputFormat {
	case "json":
		ow.WriteSingleObjectJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		fmt.Println(minorVersionUpgradeSummary(db.Name, u))
	}
}

// minorVersionUpgradeSummary is the human output of a dry run or a started
// upgrade.
func minorVersionUpgradeSummary(name string, u *civogo.UpgradeDatabaseMinorVersionResponse) string {
	if u.Started {
		return fmt.Sprintf("Upgrading %s from PostgreSQL %s to %s. The database restarts during the upgrade; follow progress with civo db show %s.", name, u.FromVersion, u.ToVersion, name)
	}
	return fmt.Sprintf("%s can be upgraded from PostgreSQL %s to %s.", name, u.FromVersion, u.ToVersion)
}

// minorVersionUpgradeWarning spells out what the upgrade does before the
// confirmation prompt.
func minorVersionUpgradeWarning(from, to string) string {
	return strings.Join([]string{
		fmt.Sprintf("Upgrading from PostgreSQL %s to %s:", from, to),
		"  - the database restarts: replicas first, then a switchover to an upgraded replica",
		"  - data, endpoint and point-in-time restore are unaffected",
		"  - there is no downgrade",
	}, "\n")
}
