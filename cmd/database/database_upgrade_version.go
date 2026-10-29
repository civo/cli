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

// errVersionUpgradeChecksFailed reports failed pre-flight checks, which are
// already printed in full.
var errVersionUpgradeChecksFailed = errors.New("pre-flight checks failed")

// runMajorVersionUpgrade upgrades one PostgreSQL database in place to the
// major given with --postgres-version and exits non-zero when it does not.
func runMajorVersionUpgrade(client *civogo.Client, db *civogo.Database) {
	err := majorVersionUpgrade(client, db, upgradePostgresVersion, upgradeDryRun, confirmMajorVersionUpgrade)
	if err == nil {
		return
	}
	if !errors.Is(err, errVersionUpgradeChecksFailed) {
		utility.Error("%s", err)
	}
	os.Exit(1)
}

// majorVersionUpgrade runs the pre-flight checks first, so a failing check is
// shown in full before anything changes. The real request is only sent when
// every check passes, it is not a dry run and confirm agrees, and it asks for
// the major the checks were run against.
func majorVersionUpgrade(client *civogo.Client, db *civogo.Database, target string, dryRun bool,
	confirm func(db *civogo.Database, plan *civogo.UpgradeDatabaseVersionResponse) bool) error {
	plan, err := client.UpgradeDatabaseVersion(db.ID, &civogo.UpgradeDatabaseVersionRequest{TargetVersion: strings.TrimSpace(target), ValidateOnly: true})
	if err != nil {
		return fmt.Errorf("failed to check the upgrade: %w", err)
	}

	failed := failedVersionUpgradeChecks(plan.Checks)
	if dryRun || len(failed) > 0 {
		writeVersionUpgradeChecks(plan.Checks)
		if common.OutputFormat == "human" {
			fmt.Println(versionUpgradeVerdict(db.Name, plan))
		}
		if len(failed) > 0 {
			return errVersionUpgradeChecksFailed
		}
		return nil
	}

	if !confirm(db, plan) {
		fmt.Println("Aborted")
		return nil
	}

	result, err := client.UpgradeDatabaseVersion(db.ID, &civogo.UpgradeDatabaseVersionRequest{TargetVersion: plan.ToVersion})
	if err != nil {
		return fmt.Errorf("failed to upgrade the database: %w", err)
	}
	writeVersionUpgrade(db, result)
	return nil
}

// confirmMajorVersionUpgrade explains what the upgrade does and asks before
// it starts, unless --yes was given. The explanation goes to stderr for
// json and custom output so stdout stays machine-readable.
func confirmMajorVersionUpgrade(db *civogo.Database, plan *civogo.UpgradeDatabaseVersionResponse) bool {
	if common.DefaultYes {
		return true
	}
	w := os.Stdout
	if common.OutputFormat != "human" {
		w = os.Stderr
	}
	fmt.Fprintln(w, majorVersionUpgradeWarning(plan.FromVersion, plan.ToVersion))
	msg := fmt.Sprintf("upgrade %s from PostgreSQL %s to %s", utility.Green(db.Name), utility.Green(plan.FromVersion), utility.Green(plan.ToVersion))
	return utility.AskForConfirm(msg) == nil
}

func writeVersionUpgrade(db *civogo.Database, result *civogo.UpgradeDatabaseVersionResponse) {
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
		fmt.Println(versionUpgradeSummary(db.Name, result))
	}
}

// versionUpgradeSummary is the human output once the upgrade was requested.
// The backup name is the one the API reports: it can carry a suffix when an
// earlier attempt kept its backup, so the CLI never guesses it.
func versionUpgradeSummary(name string, r *civogo.UpgradeDatabaseVersionResponse) string {
	if !r.Started {
		return fmt.Sprintf("The upgrade of %s from PostgreSQL %s to %s did not start; check civo db show %s.",
			utility.Green(name), utility.Green(r.FromVersion), utility.Green(r.ToVersion), name)
	}
	backup := "A backup is taken first"
	if r.BackupName != "" {
		backup = fmt.Sprintf("A backup named %s is taken first", utility.Green(r.BackupName))
	}
	return fmt.Sprintf("Upgrading %s from PostgreSQL %s to %s. %s; follow progress with civo db show %s.",
		utility.Green(name), utility.Green(r.FromVersion), utility.Green(r.ToVersion), backup, name)
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

// majorVersionUpgradeWarning spells out what the upgrade does before the
// confirmation prompt.
func majorVersionUpgradeWarning(from, to string) string {
	return strings.Join([]string{
		fmt.Sprintf("Upgrading from PostgreSQL %s to %s:", from, to),
		"  - a backup is taken first and kept; its name is shown once the upgrade starts",
		"  - the database stops for a few minutes while its data is upgraded",
		"  - replicas re-sync afterwards; until then there is no high availability or read-only endpoint",
		fmt.Sprintf("  - point-in-time restore starts again from the first backup on PostgreSQL %s", to),
		"  - there is no downgrade; to go back, restore that backup into a new database",
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
