package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/civo/cli/utility"
)

// runMajorVersionUpgrade upgrades one PostgreSQL database in place to the
// major given with --postgres-version. The pre-flight checks always run
// first, so a failing check is shown in full before anything changes.
func runMajorVersionUpgrade(client *civogo.Client, db *civogo.Database) {
	target := strings.TrimSpace(upgradePostgresVersion)
	plan, err := client.UpgradeDatabaseVersion(db.ID, &civogo.UpgradeDatabaseVersionRequest{TargetVersion: target, ValidateOnly: true})
	if err != nil {
		utility.Error("Failed to check the upgrade: %s", err)
		os.Exit(1)
	}

	failed := failedVersionUpgradeChecks(plan.Checks)
	if upgradeDryRun || len(failed) > 0 {
		writeVersionUpgradeChecks(plan.Checks)
		if common.OutputFormat == "human" {
			fmt.Println(versionUpgradeVerdict(db.Name, plan))
		}
		if len(failed) > 0 && !upgradeDryRun {
			os.Exit(1)
		}
		return
	}

	if !common.DefaultYes {
		fmt.Println(majorVersionUpgradeWarning(plan.FromVersion, plan.ToVersion, preUpgradeBackupName(plan)))
		msg := fmt.Sprintf("upgrade %s from PostgreSQL %s to %s", utility.Green(db.Name), utility.Green(plan.FromVersion), utility.Green(plan.ToVersion))
		if err := utility.AskForConfirm(msg); err != nil {
			fmt.Println("Aborted")
			return
		}
	}

	result, err := client.UpgradeDatabaseVersion(db.ID, &civogo.UpgradeDatabaseVersionRequest{TargetVersion: plan.ToVersion})
	if err != nil {
		utility.Error("Failed to upgrade the database: %s", err)
		os.Exit(1)
	}

	ow := utility.NewOutputWriterWithMap(map[string]string{
		"id":           db.ID,
		"name":         db.Name,
		"from_version": result.FromVersion,
		"to_version":   result.ToVersion,
		"backup_name":  result.BackupName,
		"started":      fmt.Sprintf("%t", result.Started),
	})
	switch common.OutputFormat {
	case "json":
		ow.WriteSingleObjectJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		fmt.Printf("Upgrading %s from PostgreSQL %s to %s. A backup named %s is taken first; follow progress with civo db show %s.\n",
			utility.Green(db.Name), utility.Green(result.FromVersion), utility.Green(result.ToVersion), utility.Green(preUpgradeBackupName(result)), db.Name)
	}
}

// writeVersionUpgradeChecks prints every pre-flight check.
func writeVersionUpgradeChecks(checks []civogo.DatabaseVersionUpgradeCheck) {
	ow := utility.NewOutputWriter()
	for _, ch := range checks {
		ow.StartLine()
		ow.AppendDataWithLabel("check", ch.Name, "Check")
		ow.AppendDataWithLabel("result", checkResult(ch.Passed), "Result")
		ow.AppendDataWithLabel("detail", ch.Message, "Detail")
	}
	switch common.OutputFormat {
	case "json":
		ow.WriteMultipleObjectsJSON(common.PrettySet)
	case "custom":
		ow.WriteCustomOutput(common.OutputFields)
	default:
		ow.WriteTable()
	}
}

func checkResult(passed bool) string {
	if passed {
		return "passed"
	}
	return "failed"
}

func failedVersionUpgradeChecks(checks []civogo.DatabaseVersionUpgradeCheck) []civogo.DatabaseVersionUpgradeCheck {
	var failed []civogo.DatabaseVersionUpgradeCheck
	for _, ch := range checks {
		if !ch.Passed {
			failed = append(failed, ch)
		}
	}
	return failed
}

// versionUpgradeVerdict is the one-line summary under the checks table.
func versionUpgradeVerdict(name string, plan *civogo.UpgradeDatabaseVersionResponse) string {
	failed := len(failedVersionUpgradeChecks(plan.Checks))
	switch failed {
	case 0:
		return fmt.Sprintf("All checks passed: %s can be upgraded from PostgreSQL %s to %s.", name, plan.FromVersion, plan.ToVersion)
	case 1:
		return fmt.Sprintf("1 check failed: fix it before upgrading %s from PostgreSQL %s to %s.", name, plan.FromVersion, plan.ToVersion)
	default:
		return fmt.Sprintf("%d checks failed: fix them before upgrading %s from PostgreSQL %s to %s.", failed, name, plan.FromVersion, plan.ToVersion)
	}
}

// preUpgradeBackupName is the backup the upgrade takes first: the name the
// API reports, or the name it uses when it reports none yet.
func preUpgradeBackupName(r *civogo.UpgradeDatabaseVersionResponse) string {
	if r.BackupName != "" {
		return r.BackupName
	}
	return "pre-upgrade-pg" + r.ToVersion
}

// majorVersionUpgradeWarning spells out what the upgrade does before the
// confirmation prompt.
func majorVersionUpgradeWarning(from, to, backup string) string {
	return strings.Join([]string{
		fmt.Sprintf("Upgrading from PostgreSQL %s to %s:", from, to),
		fmt.Sprintf("  - a backup named %s is taken first and kept", backup),
		"  - the database stops for a few minutes while its data is upgraded",
		"  - replicas re-sync afterwards; until then there is no high availability or read-only endpoint",
		fmt.Sprintf("  - point-in-time restore starts again from the first backup on PostgreSQL %s", to),
		fmt.Sprintf("  - there is no downgrade; to go back, restore %s into a new database", backup),
	}, "\n")
}

// formatVersionUpgrade summarises the latest major upgrade for civo db show.
func formatVersionUpgrade(u *civogo.DatabaseVersionUpgrade) string {
	if u == nil || u.Phase == "" {
		return ""
	}
	out := fmt.Sprintf("%s (PostgreSQL %s to %s", u.Phase, u.FromVersion, u.TargetVersion)
	if u.BackupName != "" {
		out += ", backup " + u.BackupName
	}
	out += ")"
	if u.Message != "" {
		out += ": " + u.Message
	}
	return out
}
